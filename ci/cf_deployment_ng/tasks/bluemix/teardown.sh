#!/usr/bin/env bash
set -e

: ${DEPLOYMENT_BlUEMIX_NAME:?}

source /etc/profile.d/chruby.sh

chruby 2.2.4

deployment_dir="${PWD}/deployment"
mkdir -p $deployment_dir

tar -zxvf director-artifacts/director_artifacts.tgz -C ${deployment_dir}

cp ${deployment_dir}/bosh-cli* /usr/local/bin/bosh-cli
chmod +x /usr/local/bin/bosh-cli

echo -e "\n\033[32m[INFO] Verifying the director environment.\033[0m"
cat ${deployment_dir}/director-hosts >>/etc/hosts
export BOSH_ENVIRONMENT=$(awk '{if ($2!="") print $2}' ${deployment_dir}/director-hosts)
export BOSH_CLIENT=admin
export BOSH_CLIENT_SECRET=$(bosh-cli int ${deployment_dir}/director-creds.yml --path /admin_password)
export BOSH_CA_CERT=$(bosh-cli int ${deployment_dir}/director-creds.yml --path /default_ca/ca)

echo -e "\n\033[32m[INFO] Deleting deployments and resources without inspecting the existence of deployments.\033[0m"
bosh-cli -n -d ${DEPLOYMENT_BlUEMIX_NAME} delete-deployment --force
bosh-cli -n clean-up --all

echo -e "\n\033[32m[INFO] Deleting director.\033[0m"
bosh-cli -n delete-env ${deployment_dir}/director.yml