#!/usr/bin/env bash

set -e -x

source bosh-softlayer-tools/ci/tasks/utils.sh

source /etc/profile.d/chruby.sh
chruby 2.4.4

pushd deployment
  cp -r ./.bosh_init $HOME/

  chmod +x ../bosh-init/bosh-init*

  echo "using bosh-init CLI version..."
  ../bosh-init/bosh-init* version

  echo "deleting existing BOSH Director VM..."
  ../bosh-init/bosh-init* delete director-manifest.yml
popd
