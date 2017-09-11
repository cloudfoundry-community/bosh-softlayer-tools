#!/usr/bin/env bash
set -e

source /etc/profile.d/chruby.sh

chruby 2.2.4

: ${INFRASTRUCTURE:?}
: ${SL_VM_PREFIX:?}
: ${SL_VM_DOMAIN:?}
: ${SL_USERNAME:?}
: ${SL_API_KEY:?}
: ${SL_DATACENTER:?}
: ${SL_VLAN_PUBLIC:?}
: ${SL_VLAN_PRIVATE:?}

cp bosh-cli-v2/bosh-cli-* /usr/local/bin/bosh-cli
chmod +x /usr/local/bin/bosh-cli

deployment_dir="${PWD}/deployment"
mkdir -p $deployment_dir
chmod +x bosh-cli-v2/bosh-cli*

echo -e "\n\033[32m[INFO] Using bosh-cli $(bosh-cli -v).\033[0m"
echo -e "\n\033[32m[INFO] Generating director yml.\033[0m"
cat >remove_variables.yml <<EOF
- type: remove
  path: /variables
EOF

cat >remove-health-monitor.yml <<EOF
- path: /instance_groups/name=bosh/jobs/name=health_monitor
  type: remove
EOF

echo -e "\n\033[32m[INFO] Generating manifest director.yml.\033[0m"
bosh-cli int bosh-deployment/bosh.yml \
	-o bosh-deployment/$INFRASTRUCTURE/cpi-base.yml \
	-o bosh-deployment/$INFRASTRUCTURE/dynamic-director.yml \
	-o bosh-deployment/jumpbox-user.yml \
	-o ./remove-health-monitor.yml \
	-v internal_ip=$SL_VM_PREFIX.$SL_VM_DOMAIN \
	-v dns_recursor_ip=$SL_VM_PREFIX.$SL_VM_DOMAIN \
	-v director_name=bats-director \
	-v sl_director_fqn=$SL_VM_PREFIX.$SL_VM_DOMAIN \
	-v sl_datacenter=$SL_DATACENTER \
	-v sl_vlan_public=$SL_VLAN_PUBLIC \
	-v sl_vlan_private=$SL_VLAN_PRIVATE \
	-v sl_vm_name_prefix=$SL_VM_PREFIX \
	-v sl_vm_domain=$SL_VM_DOMAIN \
	-v sl_username=$SL_USERNAME \
	-v sl_api_key=$SL_API_KEY \
	--vars-store ${deployment_dir}/director-creds.yml \
	>${deployment_dir}/director.yml

echo -e "\n\033[32m[INFO] Deploying director.\033[0m"
bosh-cli create-env \
	--state=${deployment_dir}/director-state.json \
	--vars-store ${deployment_dir}/director-creds.yml \
	${deployment_dir}/director.yml

echo -e "\n\033[32m[INFO] Deploy director successfully:\033[0m"
cat /etc/hosts | grep "$SL_VM_DOMAIN" | tee ${deployment_dir}/director-hosts

echo -e "\n\033[32m[INFO] Updating cloud-config.\033[0m"
export BOSH_ENVIRONMENT=$(awk '{if ($2!="") print $2}' ${deployment_dir}/director-hosts)
export BOSH_CLIENT=admin
export BOSH_CLIENT_SECRET=$(bosh-cli int ${deployment_dir}/director-creds.yml --path /admin_password)
export BOSH_CA_CERT=$(bosh-cli int ${deployment_dir}/director-creds.yml --path /default_ca/ca)

cat >cloud-config-plus.yml <<EOF
- type: replace
  path: /networks/name=default?
  value:
    name: default
    type: dynamic
    subnets:
    - range: ((internal_cidr))
      gateway: ((internal_gw))
      azs: [z1, z2, z3]
      dns: [8.8.8.8, 10.0.80.11, 10.0.80.12]
      cloud_properties:
        vlanIds: [((sl_vlan_public_id)), ((sl_vlan_private_id))]

- type: replace
  path: /networks/name=dynamic?
  value:
    name: dynamic
    type: dynamic
		subnets:
		- cloud_properties:
				vlanIds: [((sl_vlan_public_id)), ((sl_vlan_private_id))]
			dns: [8.8.8.8, 10.0.80.11, 10.0.80.12]
			azs: [z1]

- type: replace
  path: /vm_types?
  value:
  - name: default
    cloud_properties:
      startCpus:  4
      maxMemory:  8192
      maxNetworkSpeed: 100
      ephemeralDiskSize: 100
      hourlyBillingFlag: true
      vmNamePrefix: ((sl_vm_name_prefix))
      domain: ((sl_vm_domain))
      datacenter: ((sl_datacenter))	
  - name: minimal
    cloud_properties:
      startCpus:  1
      maxMemory:  4096
      maxNetworkSpeed: 100
      ephemeralDiskSize: 100
      hourlyBillingFlag: true
      vmNamePrefix: ((sl_vm_name_prefix))
      domain: ((sl_vm_domain))
      datacenter: ((sl_datacenter))
  - name: small
    cloud_properties:
      startCpus:  2
      maxMemory:  8192
      maxNetworkSpeed: 100
      ephemeralDiskSize: 100
      hourlyBillingFlag: true
      vmNamePrefix: ((sl_vm_name_prefix))
      domain: ((sl_vm_domain))
      datacenter: ((sl_datacenter))
  - name: small-highmem
    cloud_properties:
      startCpus:  4
      maxMemory:  32768
      maxNetworkSpeed: 100
      ephemeralDiskSize: 100
      hourlyBillingFlag: true
      vmNamePrefix: ((sl_vm_name_prefix))
      domain: ((sl_vm_domain))
      datacenter: ((sl_datacenter))
  - name: sharedcpu
    cloud_properties:
      startCpus:  1
      maxMemory:  2048
      maxNetworkSpeed: 100
      ephemeralDiskSize: 100
      hourlyBillingFlag: true
      vmNamePrefix: ((sl_vm_name_prefix))
      domain: ((sl_vm_domain))
      datacenter: ((sl_datacenter))

- type: replace
  path: /disk_types?
	value:
  - name: default
	  disk_size: 100_000
	  cloud_properties:
		  iops: 3000
		  snapshotSpace: 10
  - name: small
	  disk_size: 20_000
	  cloud_properties:
		  iops: 500
		  snapshotSpace: 5
  - name: large
	  disk_size: 250_000
	  cloud_properties:
		  iops: 3000
		  snapshotSpace: 50
  - disk_size: 5120
	  name: 5GB
  - disk_size: 10240
	  name: 10GB
  - disk_size: 100240
  name: 100GB

- type: replace
  path: /azs/name=default?
	value:
	  name: default
		cloud_properties:
			datacenter: {name: ((sl_datacenter)) }

- type: remove
  path: /networks/name=dynamic_public?

- type: remove
  path: /networks/name=dynamic_private?

- type: replace
  path: /vm_extensions?
	value:
	- name: cf-router-network-properties
	- name: cf-tcp-router-network-properties
	- name: diego-ssh-proxy-network-properties
	- name: 50GB_ephemeral_disk
		cloud_properties:
			disk: 51200
	- name: 100GB_ephemeral_disk
		cloud_properties:
			disk: 102400
EOF

bosh-cli update-cloud-config -n ./bosh-deployment/softlayer/cloud-config.yml \
	-o ./cloud-config-plus.yml \
	-v sl_datacenter=lon02 \
	-v sl_vm_name_prefix=bosh-automation-cf-test \
	-v sl_vm_domain=softlayer.com \
	-v internal_cidr=10.0.0.0/24 \
	-v internal_gw=10.0.0.1 \
	-v sl_vlan_public_id=1292653 \
	-v sl_vlan_private_id=1292651

echo -e "\n\033[32m[INFO] Final state of director deployment:\033[0m"
cat ${deployment_dir}/director-state.json

echo -e "\n\033[32m[INFO] Saving config.\033[0m"
cp bosh-cli-v2/bosh-cli* ${deployment_dir}/
pushd ${deployment_dir}
tar -zcvf /tmp/director_artifacts.tgz ./ >/dev/null 2>&1
popd
mv /tmp/director_artifacts.tgz deploy-artifacts/
