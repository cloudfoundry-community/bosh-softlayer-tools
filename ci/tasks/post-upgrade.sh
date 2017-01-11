#!/usr/bin/env bash
set -ex

dir=`dirname "$0"`
source ${dir}/utils.sh

deployment_yml=${DEPLOYMENT_YML}

bosh_cli=${BOSH_CLI}
bosh_cli_password=${BOSH_CLI_PASSWORD}
install_bosh_cli
echo "login director..."
bosh -n target ${BLUEMIX_DIRECTOR_IP}
bosh login admin admin

echo "backup deployment yml..."
sudo apt-get -y install expect
/usr/bin/env expect<<EOF
spawn scp -o StrictHostKeyChecking=no -r ./${deployment_yml} root@$bosh_cli:/root/security/
expect "*?assword:*"
exp_send "$bosh_cli_password\r"
expect eof
EOF

echo "get old versions"
old_stemcell_version=`bosh stemcells|grep bosh-softlayer-xen-ubuntu-trusty-go_agent|awk '{print $6}'|head -n 1`
old_security_version=`bosh releases|grep security-release| awk '{print $4}'`

echo yes | bosh delete release security-release $old_security_version
echo yes | bosh delete stemcell bosh-softlayer-xen-ubuntu-trusty-go_agent $old_stemcell_version