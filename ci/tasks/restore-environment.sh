#!/bin/bash
set -ex

dir=`dirname "$0"`
source ${dir}/utils.sh

bosh_cli=${BOSH_CLI}
bosh_cli_password=${BOSH_CLI_PASSWORD}
deployment_yml=${DEPLOYMENT_YML}
install_bosh_cli
echo "login director..."
bosh -n target ${BLUEMIX_DIRECTOR_IP}
bosh login admin admin

sudo apt-get -y install expect
set timeout 30
/usr/bin/env expect<<EOF
spawn scp -o StrictHostKeyChecking=no root@$bosh_cli:/root/security/${deployment_yml} ./
expect "*?assword:*"
exp_send "$bosh_cli_password\r"
expect eof
EOF

function restore (){
  print_title "RESTORING..."
  bosh deployment ${deployment_yml}
  echo "yes" | bosh deploy

  if [ $? -eq 0 ]; then
     echo "Restore successful!"
  else
     echo "Restore failed!"
  fi
}