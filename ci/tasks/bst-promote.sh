#!/usr/bin/env bash

set -e -x

version=`cat version-semver/number`

base=$( cd "$( dirname "$( dirname "$0" )")"/.. && pwd )

base_gopath=$( cd $base/../../../.. && pwd )

export GOPATH=$base/Godeps/_workspace:$base_gopath:$GOPATH

echo "GOPATH=" $GOPATH

echo -e "\n Formatting packages..."
go fmt ./...

echo -e "\nGenerating Binary: bmp..."
go build -o $base/out/bmp $base/main/bmp/bmp.go
chmod +x $base/out/bmp

echo -e "\nGenerating Binary: sl_stemcells..."
go build -o $base/out/sl_stemcells $base/main/stemcells/stemcells.go
chmod +x $base/out/sl_stemcells

echo -e "\nGenerating bosh-softlayer-tools"
tar -cvf bosh-softlayer-tools-$version.tar -C $base/out/ $base/out/*

mv bosh-softlayer-tools-$version.tar promoted/

