#!/usr/bin/env bash
set -e

source bosh-softlayer-tools/ci/tasks/utils.sh

check_param CF_API
check_param CF_USERNAME
check_param CF_PASSWORD
check_param APP_API

source bosh-softlayer-tools/ci/tasks/utils.sh
deployment_dir="${PWD}/deployment"
mkdir -p $deployment_dir

find . -maxdepth 2 -type f -name "director_artifacts*.tgz" |xargs tar -C ${deployment_dir} -xvpf

function install_cf_cli () {
  print_title "INSTALL CF CLI..."
  curl -L "https://cli.run.pivotal.io/stable?release=linux64-binary&source=github" | tar -zx
  mv cf /usr/local/bin
  echo "cf version..."
  cf --version
}

function cf_push_cpp () {
  print_title "CF PUSH APP..."
  name_server=$(cat deployment/director-hosts|awk '{print $1}')

  cat /etc/resolv.conf
  sed -i '1 i\nameserver '"${name_server}"'' /etc/resolv.conf
  app="cf-app/IICVisit.war"

  CF_TRACE=true cf api ${CF_API} --skip-ssl-validation
  CF_TRACE=true cf login -u ${CF_USERNAME} -p ${CF_PASSWORD} --skip-ssl-validation
  cf set-quota org q4GB
  cf create-space dev
  cf target -o org -s dev
  cf push IICVisit -p ${app}
  curl iicvisit.${APP_API}/GetEnv|grep "DEA IP"
  if [ $? -eq 0 ]; then
   echo "cf push app successful!"
  else
   echo "cf push app failed!"
  fi
}


install_cf_cli

cf_push_cpp