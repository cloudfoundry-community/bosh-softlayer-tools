#!/usr/bin/env bash
set -e

source /etc/profile.d/chruby.sh

chruby 2.2.4

: ${SL_VM_PREFIX:?}
: ${SL_USERNAME:?}
: ${SL_API_KEY:?}
: ${SL_DATACENTER:?}
: ${SL_VLAN_PUBLIC:?}
: ${SL_VLAN_PRIVATE:?}

cp bosh-cli-v2/bosh-cli-* /usr/local/bin/bosh-cli
chmod +x /usr/local/bin/bosh-cli

deployment_dir="${PWD}/deployment"
mkdir -p $deployment_dir
SL_VM_DOMAIN=${SL_VM_PREFIX}.softlayer.com
chmod +x bosh-cli-v2/bosh-cli*

echo -e "\n\033[32m[INFO] Using bosh-cli $(bosh-cli -v).\033[0m"
echo -e "\n\033[32m[INFO] Generating director yml.\033[0m"
cat >remove_variables.yml <<EOF
- type: remove
  path: /variables
EOF
bosh-cli int bosh-softlayer-tools/ci/templates/director-template.yml \
	-o remove_variables.yml \
	--vars-store ${deployment_dir}/credentials.yml \
	-v SL_VM_PREFIX=${SL_VM_PREFIX} \
	-v SL_VM_DOMAIN=${SL_VM_DOMAIN} \
	-v SL_USERNAME=${SL_USERNAME} \
	-v SL_API_KEY=${SL_API_KEY} \
	-v SL_DATACENTER=${SL_DATACENTER} \
	-v SL_VLAN_PUBLIC=${SL_VLAN_PUBLIC} \
	-v SL_VLAN_PRIVATE=${SL_VLAN_PRIVATE} \
	-v BOSH_RELEASE=${BOSH_RELEASE} \
	-v BOSH_RELEASE_SHA1=${BOSH_RELEASE_SHA1} \
	-v BOSH_SOFTLAYER_CPI_RELEASE=${BOSH_SOFTLAYER_CPI_RELEASE} \
	-v BOSH_SOFTLAYER_CPI_RELEASE_SHA1=${BOSH_SOFTLAYER_CPI_RELEASE_SHA1} \
	-v DIRECTOR_STEMCELL=${DIRECTOR_STEMCELL} \
	-v STEMCELL_SHA1=${STEMCELL_SHA1} \
	>${deployment_dir}/director-base.yml

echo -e "\n\033[32m[INFO] Deploying director.\033[0m"
bosh-cli create-env ${deployment_dir}/director-base.yml \
	--state=${deployment_dir}/director-deploy-state.json

echo -e "\n\033[32m[INFO] Final state of director deployment:\033[0m"

cat ${deployment_dir}/director-deploy-state.json

echo -e "\n\033[32m[INFO] Director:\033[0m"
cat /etc/hosts | grep "$SL_VM_DOMAIN" | tee ${deployment_dir}/director-hosts

echo -e "\n\033[32m[INFO] Saving config.\033[0m"
cp bosh-cli-v2/bosh-cli* ${deployment_dir}/
pushd ${deployment_dir}
  tar -zcvf /tmp/director_artifacts.tgz ./ >/dev/null 2>&1
popd
mv /tmp/director_artifacts.tgz deploy-artifacts/
