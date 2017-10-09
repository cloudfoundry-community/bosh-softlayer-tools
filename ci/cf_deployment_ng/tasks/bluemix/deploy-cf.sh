#!/usr/bin/env bash
set -e
source bosh-softlayer-tools/ci/tasks/utils.sh

check_param deploy_name
check_param data_center_name
check_param private_vlan_id
check_param public_vlan_id

deployment_dir="${PWD}/deployment"
mkdir -p $deployment_dir

tar -zxvf director-artifacts/director_artifacts.tgz -C ${deployment_dir}
cp ${deployment_dir}/bosh-cli* /usr/local/bin/bosh-cli
chmod +x /usr/local/bin/bosh-cli

echo -e "\n\033[32m[INFO] Using bosh-cli $(bosh-cli -v).\033[0m"

cat ${deployment_dir}/director-hosts >> /etc/hosts
bosh-cli -e $(cat ${deployment_dir}/director-hosts |awk '{print $2}') --ca-cert <(bosh-cli int ${deployment_dir}/director-creds.yml --path /director_ssl/ca ) alias-env bosh-test

director_password=$(bosh-cli int ${deployment_dir}/director-creds.yml --path /admin_password)
echo "Trying to login to director..."
export BOSH_CLIENT=admin
export BOSH_CLIENT_SECRET=${director_password}
bosh-cli -e bosh-test login

director_ip=$(awk '{print $1}' ${deployment_dir}/director-hosts)
director_uuid=$(grep -Po '(?<=director_id": ")[^"]*' ${deployment_dir}/director-state.json)


echo -e "\n\033[32m[INFO] Uploading stemcell.\033[0m"
bosh-cli -e bosh-test us https://s3.amazonaws.com/ng-bosh-softlayer-stemcells-bluemix-candidate-container/light-bosh-stemcell-3445.11.3-bluemix-xen-ubuntu-trusty-go_agent.tgz

# generate cf deployment yml file
bosh-cli int cf-template/${cf_template} \
							-v director_password=${director_password} \
							-v director_ip=${director_ip}\
							-v director_uuid=${director_uuid}\
							-v deploy_name=${deploy_name}\
							-v data_center_name=${data_center_name}\
							-v private_vlan_id=${private_vlan_id}\
							-v public_vlan_id=${public_vlan_id}\
							-v stemcell_version=\"${stemcell_version}\"\
							-v stemcell_location=${stemcell_location}\
							-v stemcell_name=${stemcell_name}\
						    > ${deployment_dir}/cf-deploy-base.yml

# generate diego deployment yml file
bosh-cli int diego-template/${diego_template} \
							-v director_password=${director_password} \
							-v director_ip=${director_ip}\
							-v director_uuid=${director_uuid}\
							-v deploy_name=${deploy_name}\
							-v data_center_name=${data_center_name}\
							-v private_vlan_id=${private_vlan_id}\
							-v public_vlan_id=${public_vlan_id}\
							-v stemcell_version=\"${stemcell_version}\"\
							-v stemcell_location=${stemcell_location}\
							-v stemcell_name=${stemcell_name}\
						    > ${deployment_dir}/diego-deploy-base.yml

releases_cf=$(bosh-cli int ${deployment_dir}/cf-deploy-base.yml --path /releases |grep -Po '(?<=- location: ).*')
releases_diego=$(bosh-cli int ${deployment_dir}/diego-deploy-base.yml --path /releases |grep -Po '(?<=- location: ).*')
releases=`echo -e "${releases_cf}\n${releases_diego}"`

# upload releases
while IFS= read -r line; do
bosh-cli -e bosh-test upload-release $line
done <<< "$releases"

function stemcell_exist(){
	stemcell_version=$1
	uploaded_stemcells=$(bosh-cli  -e bosh-test stemcells |awk '{print $2}'|sed s/[+*]$//)
	IFS= read -r -a stemcells<<<"$uploaded_stemcells"
	for stemcell in "$stemcells"
	do
		if [ "$stemcell_version" == "$stemcell" ];then
			return 0
		fi
	done
	return 1
}

if ! stemcell_exist ${stemcell_version}; then
	bosh-cli -e bosh-test upload-stemcell ${stemcell_location}
fi
bosh-cli -e bosh-test vms > cf-artifacts/deployed-vms
bosh-cli -n -e bosh-test -d ${deploy_name} deploy ${deployment_dir}/cf-deploy-base.yml --no-redact
bosh-cli -n -e bosh-test -d ${deploy_name}-diego deploy ${deployment_dir}/diego-deploy-base.yml --no-redact

bosh-cli -e bosh-test vms > cf-artifacts/deployed-vms
cp ${deployment_dir}/cf-deploy-base.yml  cf-artifacts/cf-deploy-base.yml
cp ${deployment_dir}/diego-deploy-base.yml  cf-artifacts/diego-deploy-base.yml

pushd cf-artifacts
   tar -zcvf  /tmp/cf_artifacts.tgz ./ >/dev/null 2>&1
popd

mv /tmp/cf_artifacts.tgz ./cf-artifacts

