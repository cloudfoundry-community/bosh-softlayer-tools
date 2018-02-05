#!/usr/bin/env bash
set -e

source /etc/profile.d/chruby.sh

chruby 2.4.2

: ${INFRASTRUCTURE:?}
: ${DEPLOY_NAME:?}
: ${SL_VM_DOMAIN:?}
: ${SL_DATACENTER:?}
: ${SL_VLAN_PUBLIC:?}
: ${SL_VLAN_PRIVATE:?}

deployment_dir="${PWD}/deployment"
mkdir -p $deployment_dir

tar -zxvf director-artifacts/director_artifacts.tgz -C ${deployment_dir}

cp ${deployment_dir}/bosh-cli* /usr/local/bin/bosh-cli
chmod +x /usr/local/bin/bosh-cli

echo -e "\n\033[32m[INFO] Using bosh-cli $(bosh-cli -v).\033[0m"

cat ${deployment_dir}/director-hosts >>/etc/hosts
export BOSH_ENVIRONMENT=$(awk '{if ($2!="") print $2}' ${deployment_dir}/director-hosts)
export BOSH_CLIENT=admin
export BOSH_CLIENT_SECRET=$(bosh-cli int ${deployment_dir}/director-creds.yml --path /admin_password)
export BOSH_CA_CERT=$(bosh-cli int ${deployment_dir}/director-creds.yml --path /default_ca/ca)

echo -e "\n\033[32m[INFO] Generating cloud-config director.yml.\033[0m"
director_ip=$(awk '{if ($1!="") print $1}' ${deployment_dir}/director-hosts)
bosh-cli int ./bosh-deployment/${INFRASTRUCTURE}/bluemix-cf-cloud-config.yml \
	-v director_ip=${director_ip} \
	-v sl_datacenter=${SL_DATACENTER} \
	-v deploy_name=${DEPLOY_NAME} \
	-v sl_vm_domain=${SL_VM_DOMAIN} \
	-v sl_vlan_public_id=${SL_VLAN_PUBLIC} \
	-v sl_vlan_private_id=${SL_VLAN_PRIVATE} \
	>${deployment_dir}/cloud-config.yml
cat ${deployment_dir}/cloud-config.yml

echo -e "\n\033[32m[INFO] Updating cloud-config.\033[0m"
bosh-cli update-cloud-config -n ${deployment_dir}/cloud-config.yml

echo -e "\n\033[32m[INFO] Uploading stemcell.\033[0m"
bosh-cli -e $(cat ${deployment_dir}/director-hosts |awk '{print $2}') --ca-cert <(bosh-cli int ${deployment_dir}/director-creds.yml --path /director_ssl/ca ) alias-env bosh-test
director_password=$(bosh-cli int ${deployment_dir}/director-creds.yml --path /admin_password)
export BOSH_CLIENT=admin
export BOSH_CLIENT_SECRET=${director_password}
bosh-cli -e bosh-test login

stemcell_location=$(bosh-cli int ${deployment_dir}/cloud-config.yml --path /stemcell_location)
bosh-cli -e bosh-test upload-stemcell ${stemcell_location} --fix

echo -e "\n\033[32m[INFO] Final state of director deployment:\033[0m"
cat ${deployment_dir}/director-state.json

echo -e "\n\033[32m[INFO] Saving director artifacts.\033[0m"

pushd ${deployment_dir}
  tar -zcvf /tmp/director_artifacts.tgz ./ >/dev/null 2>&1
popd
mv /tmp/director_artifacts.tgz deploy-artifacts/
