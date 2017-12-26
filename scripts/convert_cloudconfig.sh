#!/bin/bash

if [ $# -ne 1 ] || [ ! -f $1 ]
then
    echo "Usage: $0 cloudconfig_yaml_file"
    exit 1
fi

TARGET=$1

cp  $TARGET $TARGET.origin
echo "Original YAML file has been saved to $TARGET.origin"

# azs cloud_properties
sed -i ':a;$!{N;ba};s/Datacenter:.[ ]*Name:/datacenter:/g' $TARGET

# networks cloud_properties
sed -i ':a;$!{N;ba};s/PrimaryBackendNetworkComponent:.[ ]*NetworkVlan:.[ ]*Id: /vlan_ids: [/g' $TARGET
sed -i ':a;$!{N;ba};s/.[ ]*PrimaryNetworkComponent:.[ ]*NetworkVlan:.[ ]*Id: /, /g' $TARGET
sed -i 's/^.*vlan_ids.*/&]/' $TARGET

# vm_types cloud_properties
sed -i '/Bosh_ip:/d' $TARGET
sed -i 's/EphemeralDiskSize:/ephemeral_disk_size:/g' $TARGET
sed -i 's/HourlyBillingFlag:/hourly_billing_flag:/g' $TARGET
sed -i 's/LocalDiskFlag:/local_disk_flag:/g' $TARGET
sed -i 's/MaxMemory:/memory:/g' $TARGET
sed -i 's/StartCpus:/cpu:/g' $TARGET
sed -i 's/VmNamePrefix:/hostname_prefix:/g' $TARGET

echo "$TARGET has been converted"