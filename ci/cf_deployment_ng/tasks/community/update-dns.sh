#!/usr/bin/env bash
set -e

: ${SYSTEM_DOMAIN:?}

deployment_dir="${PWD}/deployment"
mkdir -p $deployment_dir

cf_artifacts_path=$(find ${pwd} -name "cf-artifacts*.tgz" | head -n 1)
tar -zxvf ${cf_artifacts_path} -C ${deployment_dir}

cp ${deployment_dir}/bosh-cli* /usr/local/bin/bosh-cli
chmod +x /usr/local/bin/bosh-cli

echo -e "\n\033[32m[INFO] Fetching ips of directory and haproxy.\033[0m"
director_ip=$(awk '{print $1}' ${deployment_dir}/director-hosts)
haproxy_ip=$(grep haproxy ${deployment_dir}/deployed-vms|awk '{print $4}')
if [ -z "$haproxy_ip" ]; then
	echo -e "\n\033[31m[ERROR] haproxy ips is empty.\033[0m"
	exit 1
fi

echo -e "\n\033[32m[INFO] Updating PowerDNS records.\033[0m"
cat >update_dns.sh<<EOF
cat >/tmp/update_dns.sql<<ENDSQL
DO \\\$\\\$
DECLARE new_id INTEGER;
BEGIN
    DELETE FROM domains WHERE domains.name IN  ('${SYSTEM_DOMAIN}');
    DELETE FROM records WHERE records.name IN  ('${SYSTEM_DOMAIN}', '*.${SYSTEM_DOMAIN}');
    INSERT INTO domains(name, type) VALUES('${SYSTEM_DOMAIN}', 'NATIVE');
    SELECT domains.id INTO new_id FROM domains WHERE domains.name = '${SYSTEM_DOMAIN}' ;
    INSERT INTO records(name, type, content, ttl, domain_id) VALUES('${SYSTEM_DOMAIN}','SOA','localhost hostmaster@localhost 0 10800 604800 30', 300, new_id );
    INSERT INTO records(name, type, content, ttl, domain_id) VALUES('*.${SYSTEM_DOMAIN}', 'A', '$haproxy_ip', 300, new_id);
END\\\$\\\$;
ENDSQL
/var/vcap/packages/postgres/bin/psql -U postgres -d powerdns -a -f /tmp/update_dns.sql
EOF
chmod +x update_dns.sh

bosh-cli int ${deployment_dir}/director-creds.yml --path /jumpbox_ssh/private_key > ./jumpbox.key
chmod 600 ./jumpbox.key
ssh -o StrictHostKeyChecking=no jumpbox@${director_ip} -i ./jumpbox.key 'bash -s' < update_dns.sh

echo -e "\n\033[32m[INFO] Saving cf artifacts.\033[0m"
pushd ${deployment_dir}
  tar -zcvf /tmp/cf_artifacts.tgz ./ >/dev/null 2>&1
popd
mv /tmp/cf_artifacts.tgz ./cf-artifacts
