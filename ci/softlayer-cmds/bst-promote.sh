#!/usr/bin/env bash

set -e

version=`cat version-semver/number`

base=$( cd "$( dirname "$( dirname "$0" )")"/.. && pwd )
base_gopath=$( cd $base/../../../.. && pwd )
export GOPATH=$base_gopath:$GOPATH
echo "GOPATH=" $GOPATH

echo -e "\n Formatting packages..."
go fmt ./...

echo -e "\nGenerating Binary: sl_stemcells..."
go build -o $base/out/sl_stemcells $base/main/stemcells/stemcells.go
chmod +x $base/out/sl_stemcells

cp $base/out/sl_stemcells promoted/sl_stemcells-${version}
