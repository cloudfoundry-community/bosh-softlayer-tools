#!/bin/bash


set -e -x

echo "\n[INFO] Checking enviroment variables..."
source bosh-softlayer-tools-master/ci/bluemix/tasks/utils.sh

check_param SL_USERNAME
check_param SL_API_KEY
check_param SWIFT_USERNAME
check_param SWIFT_API_KEY
check_param SWIFT_CLUSTER
check_param SWIFT_CONTAINER

export CANDIDATE_BUILD_NUMBER=$( cat version/number | sed 's/\.0$//;s/\.0$//' )

base=$( cd "$( dirname "$( dirname "$( dirname "$0" )")")" && pwd )

echo -e "\n[INFO] Get stemcell vhd filename..."
stemcell_name="bosh-stemcell-$CANDIDATE_BUILD_NUMBER-$IAAS-esxi-$OS_NAME-$OS_VERSION-go_agent"
stemcell_vhd_filename="${stemcell_name}.vhd"

echo -e "\n[INFO] Getting stemcell binary..."
export sl_stemcells=$(realpath stemcell-cmds/sl_stemcells-*)
chmod +x $sl_stemcells

echo -e "\n[INFO] Softlayer create from external source..."
IFS=':' read -ra OBJ_STORAGE_ACC_NAME <<< "$SWIFT_USERNAME"
URI="swift://${OBJ_STORAGE_ACC_NAME}@${SWIFT_CLUSTER}/${SWIFT_CONTAINER}/${stemcell_vhd_filename}"
if [ "$OS_VERSION" = "xenial" ]; then
  os_ref_code=UBUNTU_16_64
elif [ "$OS_VERSION" = "trusty" ]; then
  os_ref_code=UBUNTU_14_64
else
  echo -e "\n[INFO] Unmatched os version: ${OS_VERSION}"
  exit 1
fi

$sl_stemcells -c import-image --os-ref-code ${os_ref_code} --uri ${URI} --infrastructure "$IAAS" --public-name "light-bosh-stemcell-$CANDIDATE_BUILD_NUMBER-$IAAS-xen-ubuntu-$OS_VERSION-go_agent" \
  --public-note "Public_light_stemcell_${CANDIDATE_BUILD_NUMBER}" --public | tail -1 > "./stemcell-image/stemcell-info.json"

sleep 10
image_id=`tail -1 ./stemcell-image/stemcell-info.json | grep -Eo '[0-9]+'`
SL_USERNAME=`echo ${SL_USERNAME/@/%40}`
result=`curl -sk https://${SL_USERNAME}:${SL_API_KEY}@api.softlayer.com/rest/v3/SoftLayer_Virtual_Guest_Block_Device_Template_Group/${image_id}/setBootMode/HVM.json`
echo ${result} | grep true
