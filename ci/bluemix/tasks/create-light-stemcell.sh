#!/bin/bash

  set -e

  # outputs
  output_dir="light-stemcell"
  mkdir -p ${output_dir}

  echo -e "\n Get stemcell version..."
  stemcell_version=$(cat version/number | sed 's/\.0$//;s/\.0$//')

  base=$( cd "$( dirname "$( dirname "$( dirname "$0" )")")" && pwd )
  base_gopath=$( cd $base/../../../.. && pwd )

  export GOPATH=$base_gopath:$GOPATH

  pushd stemcell
    echo -e "\n Unpacking stemcell raw tgz..."
    tar -zxvf *.tgz
    ls -la
    cp dev_tools_file_list.txt stemcell_dpkg_l.txt "${base}"
  popd

  echo -e "\n Creating stemcell binary..."
  cd "${base}"
  go build -o out/sl_stemcells main/stemcells/stemcells.go

  echo -e "\n Softlayer creating light stemcell..."
  out/sl_stemcells -c light-stemcell --version ${stemcell_version} --infrastructure "$IAAS" --stemcell-info-filename "${base_gopath}/../stemcell-info/stemcell-info.json"

  stemcell_filename=`ls light*.tgz`
  tar -zxvf ${stemcell_filename}
  tar -zcvf ${stemcell_filename} image stemcell.MF dev_tools_file_list.txt stemcell_dpkg_l.txt
  cp *.tgz "${base_gopath}/../${output_dir}/"

  checksum="$(sha1sum "${base_gopath}/../${output_dir}/${stemcell_filename}" | awk '{print $1}')"
  echo "$stemcell_filename sha1=$checksum"
