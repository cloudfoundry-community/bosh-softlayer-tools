#!/usr/bin/env bash
set -e

: ${DEPLOYMENT_NAME:?}
: ${CF_PREFIX:?}
: ${CF_DOMAIN:?}

deployment_dir="${PWD}/deployment"
mkdir -p $deployment_dir

tar -zxvf director-artifacts/director_artifacts.tgz -C ${deployment_dir}

cp ${deployment_dir}/bosh-cli* /usr/local/bin/bosh-cli
chmod +x /usr/local/bin/bosh-cli

echo -e "\n\033[32m[INFO] Targeting the BOSH Director.\033[0m"
cat ${deployment_dir}/director-hosts >>/etc/hosts
bosh-cli -e $(cat ${deployment_dir}/director-hosts | awk '{print $2}') --ca-cert <(bosh-cli int ${deployment_dir}/director-creds.yml --path /default_ca/ca) alias-env automation-cf

echo -e "\n\033[32m[INFO] Login to the director.\033[0m"
export BOSH_ENVIRONMENT=$(awk '{if ($2!="") print $2}' ${deployment_dir}/director-hosts)
export BOSH_CLIENT=admin
export BOSH_CLIENT_SECRET=$(bosh-cli int ${deployment_dir}/director-creds.yml --path /admin_password)
export BOSH_CA_CERT=$(bosh-cli int ${deployment_dir}/director-creds.yml --path /default_ca/ca)
bosh-cli -e automation-cf login

echo -e "\n\033[32m[INFO] Uploading stemcell.\033[0m"
bosh-cli us https://s3.amazonaws.com/bosh-softlayer-stemcells-candidate-container/light-bosh-stemcell-3421.11-softlayer-xen-ubuntu-trusty-go_agent.tgz

echo -e "\n\033[32m[INFO] Generating cf manifest.\033[0m"
cat >stemcell.yml <<EOF
- type: replace
  path: /stemcells/alias=default?
  value:
    alias: default
    os: ubuntu-trusty
    version: latest
EOF

cat >webdav-blobstore.yml <<EOF
- type: replace
  path: /instance_groups/name=blobstore/jobs/name=blobstore/properties/blobstore/internal_access_rules?
  value:
    - "allow 10.0.0.0/8;"
    - "allow 172.16.0.0/12;"
    - "allow 192.168.0.0/16;"
    - "allow 169.50.0.0/16;"
EOF

bosh-cli vms >cf-artifacts/deployed-vms
bosh-cli int cf-deployment/cf-deployment.yml \
	--vars-store ${DEPLOYMENT_NAME}/cf-creds.yml \
	-o stemcell.yml \
	-o cf-deployment/operations/rename-deployment.yml \
	-v deployment_name=${DEPLOYMENT_NAME} \
	-v system_domain=${CF_PREFIX}.${CF_DOMAIN} \
	> ${deployment_dir}/cf.yml
cat ${deployment_dir}/cf.yml

echo -e "\n\033[32m[INFO] Deploying CF.\033[0m"
bosh-cli -d ${DEPLOYMENT_NAME} -n deploy ${deployment_dir}/cf.yml

bosh-cli vms >cf-artifacts/deployed-vms
cp ${deployment_dir}/cf.yml cf-artifacts/cf.yml
cp ${deployment_dir}/cf-creds.yml cf-artifacts/cf-creds.yml

pushd cf-artifacts
  tar -zcvf /tmp/cf_artifacts.tgz ./ >/dev/null 2>&1
popd

mv /tmp/cf_artifacts.tgz ./cf-artifacts
