#!/usr/bin/env bash

set -e -x

trap clean_vagrant EXIT

set_up_vagrant_private_key() {
  key_path=$(mktemp -d /tmp/ssh_key.XXXXXXXXXX)/value
  echo "$BOSH_PRIVATE_KEY" > $key_path
  chmod 600 $key_path
  export BOSH_VAGRANT_KEY_PATH=$key_path
  eval `ssh-agent`
  ssh-add $key_path
}

clean_vagrant() {
  vagrant destroy remote -f || true
}

get_ip_from_vagrant_ssh_config() {
  config=$(vagrant ssh-config remote)
  echo $(echo "$config" | grep HostName | awk '{print $2}')
}

build_num=$(cat stemcell-version/number | sed 's/\.0$//;s/\.0$//')

pushd bosh-src

bundle

cd bosh-stemcell

set_up_vagrant_private_key

vagrant up remote --provider=aws

vagrant ssh -c "
  cd /bosh
  bundle
  export CANDIDATE_BUILD_NUMBER=$build_num
  export BOSH_MICRO_ENABLED=no
  bundle exec rake stemcell:build[$IAAS,$HYPERVISOR,$OS_NAME,$OS_VERSION,go,bosh-os-images,bosh-$OS_NAME-$OS_VERSION-os-image.tgz]
" remote

builder_ip=$(get_ip_from_vagrant_ssh_config)

popd

ssh ubuntu@${builder_ip} "sudo chroot /mnt/stemcells/softlayer/esxi/ubuntu/work/work/chroot touch /bosh-stemcell-${build_num}-softlayer-esxi-ubuntu-trusty-raw.tgz; sudo chroot /mnt/stemcells/softlayer/esxi/ubuntu/work/work/chroot tar zcvf /bosh-stemcell-${build_num}-softlayer-esxi-ubuntu-trusty-raw.tgz . --exclude proc --exclude dev; sudo mv /mnt/stemcells/softlayer/esxi/ubuntu/work/work/chroot/*.tgz /bosh/tmp"
scp ubuntu@${builder_ip}:/bosh/tmp/*-raw.tgz build/

