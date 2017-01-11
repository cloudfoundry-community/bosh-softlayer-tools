#!/usr/bin/env bash

set -ex

dir=`dirname "$0"`
source ${dir}/utils.sh

function get_old_new_versions (){
  print_title "GET OLD AND NEW VERSIONS OF STEMCELL/SECURITY RELEASE..."
  old_stemcell_version=`bosh stemcells|grep bosh-softlayer-xen-ubuntu-trusty-go_agent|awk '{print $6}'|sed 's/\*//g'`
  new_stemcell_version=`ls ./stemcell|grep light-bosh-stemcell| cut -d "-" -f 4`

  old_security_version=`bosh releases|grep security-release| awk '{print $4}'|sed 's/\*//g'`
  curl http://10.106.192.96/releases/security-release/|grep href|cut -d '"' -f 2|sed 's/\///g' > tmp.file
  new_security_version=`grep ^v[0-9]\.[1-9]\-[0-9]\*\$ tmp.file | tail -n 1`
}

function upload_releases (){
  print_title "UPLOAD STEMCELL AND SECURITY RELEASE..."
  bosh upload stemcell ./stemcell/light-bosh-stemcell-*.tgz --skip-if-exists
  mkdir security-release
  wget http://10.106.192.96/releases/security-release/${new_security_version}/security-release.tgz -P ./security-release/
  bosh upload release ./security-release/security-release.tgz --skip-if-exists
}

function update_deployment_yml (){
  print_title "UPDATE DEPLOYMENT YML..."
  sudo apt-get -y install expect
  set timeout 30
  /usr/bin/env expect<<EOF
  spawn scp -o StrictHostKeyChecking=no root@$bosh_cli:/root/security/${deployment_yml} ./
  expect "*?assword:*"
  exp_send "$bosh_cli_password\r"
  expect eof
EOF

  sed -i '/stemcell_version=/ s/'"$old_stemcell_version"'/'"$new_stemcell_version"'/' ./${deployment_yml}
  sleep 3
  sed -i '/security-release.tgz/ n;N;N;s/'"$old_security_version"'/'"$new_security_version"'/' ./${deployment_yml}
}

function bosh_deploy (){
  print_title "DEPLOY..."
  bosh deployment ${deployment_yml}
  echo "yes" | bosh deploy

  if [ $? -eq 0 ]; then
     echo "Deploy successful!"
  else
     echo "Deploy failed!"
  fi
}

bosh_cli=${BOSH_CLI}
bosh_cli_password=${BOSH_CLI_PASSWORD}
deployment_yml=${DEPLOYMENT_YML}

install_bosh_cli
echo "login director..."
bosh -n target ${BLUEMIX_DIRECTOR_IP}
bosh login admin admin

get_old_new_versions

upload_releases

update_deployment_yml

bosh_deploy