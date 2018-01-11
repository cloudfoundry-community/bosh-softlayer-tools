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

tar -zxvf director-artifacts/director_artifacts.tgz -C ${deployment_dir}

SL_VM_DOMAIN=${SL_VM_PREFIX}.softlayer.com

chmod +x bosh-cli-v2/bosh-cli* 

cat ${deployment_dir}/director-hosts >> /etc/hosts

director_ip=$(awk '{print $1}' ${deployment_dir}/director-hosts)

cat ${deployment_dir}/director-deploy-state.json |jq '.+{"current_ip": $director_ip}' --arg director_ip ${director_ip} |tee ${deployment_dir}/director-deploy-state.json

  function finish {
    echo "Final state of director deployment:"
    echo "====================================================================="
    cat ${deployment_dir}/director-deploy-state.json
    echo "====================================================================="git pi
    echo "Saving config..."
    pushd ${deployment_dir}
    tar -zcvf  /tmp/director_artifacts_updated.tgz ./ >/dev/null 2>&1
    popd
    mv /tmp/director_artifacts_updated.tgz deploy-artifacts-updated/
  }

trap finish EXIT

echo "Using bosh-cli $(bosh-cli-v2/bosh-cli* -v)"

bosh-cli-v2/bosh-cli* -e $(cat ${deployment_dir}/director-hosts |awk '{print $2}') --ca-cert <(bosh-cli-v2/bosh-cli* int ${deployment_dir}/credentials.yml --path /DIRECTOR_SSL/ca ) alias-env bosh-test
director_password=$(bosh-cli-v2/bosh-cli* int ${deployment_dir}/credentials.yml --path /DI_ADMIN_PASSWORD)
print_title "Trying to login to director..."
export BOSH_CLIENT=admin
export BOSH_CLIENT_SECRET=${director_password}
bosh-cli-v2/bosh-cli* -e bosh-test login
print_title "Ensure director is base version..."
#bosh-cli-v2/bosh-cli* create-env  ${deployment_dir}/director-base.yml \
#                --state=${deployment_dir}/director-deploy-state.json

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
                      > director-update-temp.yml

cat >remove_variables.yml <<EOF
- type: remove
  path: /variables
EOF

bosh-cli-v2/bosh-cli* int director-update-temp.yml -o remove_variables.yml > ${deployment_dir}/director-update.yml
print_title "Updating director..."

bosh-cli-v2/bosh-cli* create-env  ${deployment_dir}/director-update.yml \
                      --state=${deployment_dir}/director-deploy-state.json