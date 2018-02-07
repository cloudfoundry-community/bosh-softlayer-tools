#!/bin/bash

set -e -x

# outputs
output_dir="light-stemcell"
mkdir -p ${output_dir}

echo -e "\n[INFO] Get stemcell version..."
stemcell_version=$(cat version/number | sed 's/\.0$//;s/\.0$//')

base=$(pwd)

pushd stemcell
  echo -e "\n Unpacking stemcell raw tgz..."
  tar -zxvf *.tgz
  ls -la
  cp dev_tools_file_list.txt stemcell_dpkg_l.txt "${base}"
popd

echo -e "\n[INFO] Getting stemcell binary..."
export sl_stemcells=$(realpath stemcell-cmds/sl_stemcells-*)
chmod +x $sl_stemcells

echo -e "\n[INFO] Softlayer creating light stemcell..."
$sl_stemcells -c light-stemcell --version ${stemcell_version} --os-name "ubuntu-$OS_VERSION" --infrastructure "$IAAS" --stemcell-info-filename "./stemcell-info/stemcell-info.json"

echo -e "\n[INFO] Repacking light stemcell with info files..."
stemcell_filename=`ls light*.tgz`
tar -zxvf ${stemcell_filename}
tar -zcvf ${stemcell_filename} image stemcell.MF dev_tools_file_list.txt stemcell_dpkg_l.txt
cp *.tgz "./${output_dir}/"

echo -e "\n[INFO] Printing stemcell sha1 value..."
checksum="$(sha1sum "./${output_dir}/${stemcell_filename}" | awk '{print $1}')"
echo "$stemcell_filename sha1=$checksum"
