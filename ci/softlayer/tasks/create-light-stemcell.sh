#!/bin/bash

set -e -x

if [ "$STEMCELL_FORMATS" == 'replace-me' ]; then
  echo -e "\n[INFO] Set stemcell-formats for backward compatiblity."
  STEMCELL_FORMATS="softlayer-light-legacy"
fi
echo -e "\n[INFO] stemcell-formats: ${STEMCELL_FORMATS}"

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
  cp dev_tools_file_list.txt "${base}"
  if [ -f stemcell_dpkg_l.txt ]; then
     cp stemcell_dpkg_l.txt "${base}"
  fi
  if [ -f packages.txt ]; then
     cp packages.txt "${base}"
  fi
popd

echo -e "\n[INFO] Getting stemcell binary..."
export sl_stemcells=$(realpath stemcell-cmds/sl_stemcells-*)
chmod +x $sl_stemcells

echo -e "\n[INFO] Softlayer creating light stemcell..."
$sl_stemcells -c light-stemcell --version ${stemcell_version} --os-name "ubuntu-$OS_VERSION" --infrastructure "$IAAS" --stemcell-info-filename "./stemcell-info/stemcell-info.json" \
  --stemcell-formats "$STEMCELL_FORMATS"

echo -e "\n[INFO] Repacking light stemcell with info files..."
stemcell_filename=`ls light*.tgz`
tar -zxvf ${stemcell_filename}
if [ -f stemcell_dpkg_l.txt ]; then
    tar -zcvf ${stemcell_filename} image stemcell.MF dev_tools_file_list.txt stemcell_dpkg_l.txt
fi
if [ -f packages.txt ]; then
    tar -zcvf ${stemcell_filename} image stemcell.MF dev_tools_file_list.txt packages.txt
fi
cp *.tgz "./${output_dir}/"

echo -e "\n[INFO] Printing stemcell sha1 value..."
checksum="$(sha1sum "./${output_dir}/${stemcell_filename}" | awk '{print $1}')"
echo "$stemcell_filename sha1=$checksum"
