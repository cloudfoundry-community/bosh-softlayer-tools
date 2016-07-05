#!/usr/bin/env bash

set -e -x

build_num=$(cat stemcell-version/number | sed 's/\.0$//;s/\.0$//')

wget -P build/ http://10.113.109.244/mattcui/bosh-stemcell-${build_num}-softlayer-esxi-ubuntu-trusty-go_agent-raw.tgz

