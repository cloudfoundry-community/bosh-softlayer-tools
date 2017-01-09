#!/usr/bin/env bash
set -e

dir=`dirname "$0"`
source ${dir}/utils.sh

function install_cf_cli () {
  print_title "INSTALL CF CLI..."
  curl -L "https://cli.run.pivotal.io/stable?release=linux64-binary&source=github" | tar -zx
  mv cf /usr/local/bin
  echo "cf version..."
  cf --version
}

function cf_push_cpp () {
  print_title "CF PUSH APP..."
  name_server={NAME_SERVER}
  sed -i '1 i\nameserver '"${name_server}"'' /etc/resolv.conf
  sudo apt-get -y install expect
  app="IICVisit.war"
  set timeout 30
  /usr/bin/env expect<<EOF
  spawn scp -o StrictHostKeyChecking=no root@$bosh_cli:/root/security/${app} ./
  expect "*?assword:*"
  exp_send "$bosh_cli_password\r"
  expect eof
EOF

  CF_TRACE=true cf api ${CF-API}
  CF_TRACE=true cf login -u ${CF-USERNAME} -p ${CF-PASSWORD}
  base=`dirname "$0"`
  cf push IICVisit -p ${base}/${app}
  curl iicvisit.${APP-API}/GetEnv|grep "DEA IP"
  if [ $? -eq 0 ]; then
   echo "cf push app successful!"
  else
   echo "cf push app failed!"
  fi
}

install_cf_cli

cf_push_cpp