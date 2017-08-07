#!/usr/bin/env bash

set -e -x

set_up_bm_private_key() {
  key_path=$(mktemp -d /tmp/ssh_key_bm_stemcell.XXXXXXXXXX)/value
  echo "${bm_stemcell_private_key}" > $key_path
  chmod 600 $key_path
  eval `ssh-agent`
  ssh-add $key_path
}

export build_num=$(cat stemcell-version/number | sed 's/\.0$//;s/\.0$//')

set_up_bm_private_key

pushd stemcell-ubuntu-trusty-raw
scp -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no *.tgz root@${bm_sjc01}:/tmp
popd

pushd bosh-softlayer-tools
sed -i "s/%%bm_root_password%%/${bm_root_password}/g" scripts/make-fsa-on-baremetal.sh
scp -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no scripts/make-fsa-on-baremetal.sh root@${bm_sjc01}:/tmp
popd

ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@${bm_sjc01} "mkdir -p /tmp/bm_stemcell; export build_num=${build_num}; mv /tmp/make-fsa-on-baremetal.sh /tmp/bm_stemcell/; chmod +x /tmp/bm_stemcell/make-fsa-on-baremetal.sh; /tmp/bm_stemcell/make-fsa-on-baremetal.sh"
scp -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@${bm_sjc01}:/fsarchiver/*.fsa build/

