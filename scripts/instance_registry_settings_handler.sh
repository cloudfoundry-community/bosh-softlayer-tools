#!/bin/bash
set -e

usage() {
    echo "Usage: $0 -ip director_ip -p director_user_password -id instance_id"
    echo "       [-du director_username] [-d database] [-u db_username] [-f dump_file_location] [-a action]"
    echo -e "\nRequired:"
    echo "  -ip <director_ip> Specify the director ip"
    echo "  -p  <director_user_password> Specify the user's password of the director"
    echo "  -id <instance_id> Specify the instance's id which you want to dump"
    echo -e "\nOptional:"
    echo "  -a  <action> Specify the action(fetch/update) of handling registry settings,"
    echo "      default action is 'fetch': fetch instance's settings"
    echo "      Update instance's settings needs restart agent service(sv restart agent) after success"
    echo "  -du <director_username> Specify the username of the director,"
    echo "      default username is 'vcap'"
    echo "  -d  <database> Specify the database where registry is in,"
    echo "      default database is 'bosh'"
    echo "  -u  <db_username> Specify the database username,"
    echo "      default db_username is 'postgres'"
    echo "  -f  <dump_file_location> Specify the location of dump file,"
    echo "      default dump_file_location is current directory: $(pwd)"
    echo -e "\nExample:"
    echo "Fetch instance settings from director: $0 -ip 10.10.10.47 -p password -id 12345678"
    echo "Update instance settings to director: $0 -a update -ip 10.165.227.47 -p password -id 12345678 -f 2017-09-26_0/data/instance_12345678_settings.8_29.json"
}

parse_command_line()
{
    while [ $# -gt 0 ]; do
        CPARM=$1;
        shift
        case $CPARM in
            -ip)
                DIRECTOR_IP=$1; shift
            ;;
            -p)
                USER_PASSWORD=$1; shift
            ;;
            -id)
                INSTANCE_ID=$1; shift
            ;;
            -du)
                DIRECTOR_USERNAME=$1; shift
            ;;
            -d)
                DATABASE=$1; shift
            ;;
            -u)
                USERNAME=$1; shift
            ;;
            -f)
                FILE_LOCATION=$1; shift
            ;;
            -a)
                ACTION=$1; shift
            ;;
            -h)
                usage; exit 0
            ;;
            *)
                echo "ERROR: unrecognized parameter $CPARM"; usage; exit 1
            ;;
        esac
    done
    # check parameters provided
    if [ -z "$DIRECTOR_IP" ]; then
        echo "ERROR: director ip has not been specified."; usage; exit 1
    fi
    if [ -z "$USER_PASSWORD" ]; then
        echo "ERROR: director user password has not been specified."; usage; exit 1
    fi
    if [ -z "$INSTANCE_ID" ]; then
        echo "ERROR: instance id has not been specified."; usage; exit 1
    fi
    if [ -z "$DIRECTOR_USERNAME" ]; then
        DIRECTOR_USERNAME="vcap"
    fi
    if [ -z "$DATABASE" ]; then
        DATABASE="bosh"
    fi
    if [ -z "$USERNAME" ]; then
        USERNAME="postgres"
    fi
    if [ -z "$FILE_LOCATION" ]; then
        FILE_LOCATION=$(pwd)
    fi
    if [ -z "$ACTION" ]; then
        ACTION="fetch"
    fi
}

# Args: $1: director ip, $2: director username, $3 director user password, $4: database, $5 dababase username, $6 instance_id, $7 file_location
fetch_registry_settings() {
    dump_file=instance_$6_settings.$(date -d now "+%F_%H_%M").json
    # Generate script
    cat >fetch_registry_settings.sh<<EOF
psql_path=\$(dirname \`pgrep -a postg |awk '\$0 ~ "/var/vcap/" {print \$2;}'\`)

if [[ -n "\$psql_path" ]]; then
     echo "postgres path: \${psql_path}"
     psql=\${psql_path}/psql
else
    echo "Postgres path cannot be found. Exiting"
    exit 2
fi

postgres_folder=\$(ls -dt /var/vcap/store/postgres* | head -1)
echo "postgres folder: \${postgres_folder}"
cd \${postgres_folder}

\$psql -U $5 -d $4 -a -c "COPY(SELECT settings FROM registry_instances WHERE instance_id = '$6') TO STDOUT" > ./${dump_file}
EOF
    chmod +x ./fetch_registry_settings.sh
    sshpass -p $3 ssh -o StrictHostKeychecking=no $2@$1 'bash -s' < "./fetch_registry_settings.sh"
    sshpass -p $3 scp $2@$1:/var/vcap/store/postgres*/${dump_file} $7/${dump_file}
    # delete first command line
    sed -i '1d' $7/${dump_file}
    # new softlayer cpi can use jumpbox user to ssh access
    # ssh -o StrictHostKeyChecking=no jumpbox@${director_ip} -i ./jumpbox.key 'bash -s' < fetch_registry_settings.sh
}

# Args: $1: director ip, $2: director username, $3 director user password, $4: database, $5 dababase username, $6 instance_id, $7 file_location
update_registry_settings() {
    if [ ! -f $7 ]; then
        echo "ERROR: data file is not exist."; usage; exit 1
    fi
    
    # Get the last line of sql file
    modifiedSettings=$(awk 'END{print}' $7)
    
    # Generate script
    cat >update_registry_settings.sh<<EOF
psql_path=\$(dirname \`pgrep -a postg |awk '\$0 ~ "/var/vcap/" {print \$2;}'\`)

if [[ -n "\$psql_path" ]]; then
     echo "postgres path: \${psql_path}"
     psql=\${psql_path}/psql
else
    echo "Postgres path cannot be found. Exiting"
    exit 2
fi

cat >/tmp/update_registry_settings.sql<<ENDSQL
UPDATE registry_instances
SET settings = '${modifiedSettings}'
WHERE instance_id = '$6';
ENDSQL

\$psql -U $5 -d $4 -a -f /tmp/update_registry_settings.sql
EOF
    chmod +x ./update_registry_settings.sh
    sshpass -p $3 ssh -o StrictHostKeychecking=no $2@$1 'bash -s' < "./update_registry_settings.sh"
}

if [ $# -lt 1 ]; then
    usage
    exit 1
fi

parse_command_line "$@"

if [ ${ACTION} = "fetch" ]; then
    fetch_registry_settings ${DIRECTOR_IP} ${DIRECTOR_USERNAME} ${USER_PASSWORD} ${DATABASE} ${USERNAME} ${INSTANCE_ID} ${FILE_LOCATION}
elif [ ${ACTION} = "update" ]; then
    update_registry_settings ${DIRECTOR_IP} ${DIRECTOR_USERNAME} ${USER_PASSWORD} ${DATABASE} ${USERNAME} ${INSTANCE_ID} ${FILE_LOCATION}
fi
