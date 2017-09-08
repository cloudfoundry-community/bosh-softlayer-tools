#!/usr/bin/env bash
set -e

: ${DEPLOYMENT_NAME:?}
: ${SL_DATACENTER:?}
: ${SL_VLAN_PUBLIC:?}
: ${SL_VLAN_PRIVATE:?}
: ${CF_PREFIX:?}

deployment_dir="${PWD}/deployment"
mkdir -p $deployment_dir

tar -zxvf director-artifacts/director_artifacts.tgz -C ${deployment_dir}

cp ${deployment_dir}/bosh-cli* /usr/local/bin/bosh-cli
chmod +x /usr/local/bin/bosh-cli

echo -e "\n\033[32m[INFO] Targeting the BOSH Director.\033[0m"
cat ${deployment_dir}/director-hosts >>/etc/hosts
bosh-cli -e $(cat ${deployment_dir}/director-hosts | awk '{print $2}') --ca-cert <(bosh-cli int ${deployment_dir}/director-creds.yml --path /default_ca/ca) alias-env automation-cf

echo -e "\n\033[32m[INFO] Login to the director.\033[0m"
export BOSH_CLIENT=admin
export BOSH_CLIENT_SECRET=$(bosh-cli int ${deployment_dir}/director-creds.yml --path /admin_password)
export BOSH_CA_CERT=$(bosh-cli int ${deployment_dir}/director-creds.yml --path /default_ca/ca)
bosh-cli -e automation-cf login

director_ip=$(awk '{print $1}' ${deployment_dir}/director-hosts)
director_uuid=$(grep -Po '(?<=director_id": ")[^"]*' ${deployment_dir}/director-state.json)

echo -e "\n\033[32m[INFO] Deploying CF.\033[0m"
cat >stemcell.yml <<EOF
- type: replace
  path: /stemcells/alias=default?
  value:
    - alias: default
      os: ubuntu-trusty
      version: "latest"
EOF

bosh-cli -e automation-cf vms >cf-artifacts/deployed-vms
bosh-cli -e automation-cf -d ${DEPLOYMENT_NAME} deploy cf-deployment/cf-deployment.yml \
	--vars-store env-repo/deployment-vars.yml \
	-o stemcell.yml \
	-o cf-deployment/operations/rename-deployment.yml \
	-v deployment_name=${DEPLOYMENT_NAME} \
	-v system_domain=${CF_PREFIX}

bosh-cli -e automation-cf vms >cf-artifacts/deployed-vms
cp ${deployment_dir}/cf-deploy-base.yml cf-artifacts/cf-deploy-base.yml
cp ${deployment_dir}/diego-deploy-base.yml cf-artifacts/diego-deploy-base.yml

pushd cf-artifacts
tar -zcvf /tmp/cf_artifacts.tgz ./ >/dev/null 2>&1
popd

mv /tmp/cf_artifacts.tgz ./cf-artifacts
