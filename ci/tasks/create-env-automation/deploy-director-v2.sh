#!/usr/bin/env bash
set -e

source bosh-softlayer-tools/ci/tasks/utils.sh
source /etc/profile.d/chruby.sh

chruby 2.4.2

check_param SL_VM_PREFIX
check_param SL_USERNAME
check_param SL_API_KEY
check_param SL_DATACENTER
check_param SL_VLAN_PUBLIC
check_param SL_VLAN_PRIVATE


deployment_dir="${PWD}/deployment"
mkdir -p $deployment_dir
SL_VM_DOMAIN=${SL_VM_PREFIX}.softlayer.com
chmod +x bosh-cli-v2/bosh-cli* 

  function finish {
    echo "Final state of director deployment:"
    echo "====================================================================="
    cat ${deployment_dir}/director-deploy-state.json
    echo "====================================================================="
    echo "Director:"
    echo "====================================================================="
    cat /etc/hosts | grep "$SL_VM_DOMAIN" | tee ${deployment_dir}/director-hosts
    echo "====================================================================="
    echo "Saving config..."
    cp bosh-cli-v2/bosh-cli* ${deployment_dir}/
    pushd ${deployment_dir}
    tar -zcvf  /tmp/director_artifacts.tgz ./ >/dev/null 2>&1
    popd
    mv /tmp/director_artifacts.tgz deploy-artifacts/
  }

trap finish EXIT

echo "Using bosh-cli $(bosh-cli-v2/bosh-cli* -v)"
print_title "Generating director yml..."
bosh-cli-v2/bosh-cli* int bosh-softlayer-tools/ci/templates/director-template.yml \
                      --vars-store ${deployment_dir}/credentials.yml \
                      -v SL_VM_PREFIX=${SL_VM_PREFIX} \
                      -v SL_VM_DOMAIN=${SL_VM_DOMAIN} \
                      -v SL_USERNAME=${SL_USERNAME} \
                      -v SL_API_KEY=${SL_API_KEY} \
                      -v SL_DATACENTER=${SL_DATACENTER} \
                      -v SL_VLAN_PUBLIC=${SL_VLAN_PUBLIC} \
                      -v SL_VLAN_PRIVATE=${SL_VLAN_PRIVATE} \
                      -v BOSH_RELEASE=${BOSH_RELEASE}\
                      -v BOSH_RELEASE_SHA1=${BOSH_RELEASE_SHA1}\
                      -v BOSH_SOFTLAYER_CPI_RELEASE=${BOSH_SOFTLAYER_CPI_RELEASE}\
                      -v BOSH_SOFTLAYER_CPI_RELEASE_SHA1=${BOSH_SOFTLAYER_CPI_RELEASE_SHA1}\
                      -v DIRECTOR_STEMCELL=${DIRECTOR_STEMCELL}\
                      -v STEMCELL_SHA1=${STEMCELL_SHA1}\
                      > director-base-temp.yml

cat >remove_variables.yml <<EOF
- type: remove
  path: /variables
EOF

bosh-cli-v2/bosh-cli* int director-base-temp.yml -o remove_variables.yml > ${deployment_dir}/director-base.yml

print_title "Deploying director..."

bosh-cli-v2/bosh-cli* create-env  ${deployment_dir}/director-base.yml \
                      --state=${deployment_dir}/director-deploy-state.json

