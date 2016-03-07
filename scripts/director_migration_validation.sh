#!/bin/bash

usage()
{
        echo "***********************************************************************"
        echo "* Validate vm_cid disk_cid and agent_id before migration               "    
        echo "*                                                                      "
        echo "* Usage: $0 <deployment_yml> <vm_pool_db> <user> <key>                 "
        echo "*                                                                      "
        echo "***********************************************************************"
}

if [ $# -eq 4 ]; then
        bosh_deployment_yml=$1
        vm_pool_db=$2
        user=$3
	key=$4
else
        usage
        exit 1
fi

Pre-check()
{
	if [ ! -s $bosh_deployment_yml ]; then
		echo "$bosh_deployment_yml doesn't exist, please check input file!"
		exit 1	
	elif [ ! -s $vm_pool_db ]; then
		echo "$vm_pool_db doesn't exist, please check input file!"
		exit 1
	fi

	if [[ $user == *"@"* ]]; then
		user=${user//@/%40}
	fi

	curl_cmd="curl -s https://$user:$key@api.softlayer.com/rest/v3/"

}

Get_ip()
{
	if [ `bosh target | grep "Current target is https:" | wc -l` == "1" ]; then
		director_ip=`bosh target | cut -d: -f 2 | awk -F// '{print $NF}'`
		echo "director server ip is $director_ip"
	else
		echo "Director server is not specified!"
		exit 1
	fi
}


Check_vm()
{
	if [ `grep ":vm_cid:" $bosh_deployment_yml | wc -l` != "1" ]; then
		echo "vm_cid existence validation: failed, it seems configuration vm_cid incorrect, please check $bosh_deployment_yml"
		return 1
	else
		vm_cid=`grep ":vm_cid:" $bosh_deployment_yml | awk -F\' '{print $2}'`	
	fi

	actual_cid=`${curl_cmd}SoftLayer_Virtual_Guest/findByIpAddress/$director_ip.json | grep -Po '(?<="id":)[^,]*'`
	if [ "$actual_cid" == "$vm_cid" ]; then
		echo "vm_cid existence validation: pass"
	else
		echo "vm_cid existence validation: failed  The actual vm_cid should be $actual_cid"
		vm_cid=""
	fi
}

Check_disk()
{
	if [ `grep ":disk_cid:" $bosh_deployment_yml | wc -l` != "1" ]; then
		echo "disk_cid existence validation: failed, it seems configuration disk_cid incorrect, please check $bosh_deployment_yml"
		return 1
	else
		disk_cid=`grep ":disk_cid:" $bosh_deployment_yml | awk -F\' '{print $2}'`
	fi

	if [ `${curl_cmd}SoftLayer_Network_Storage_Iscsi/getObject/$disk_cid.json | grep "\"id\":$disk_cid" |wc -l` == "1" ]; then
		echo "disk_cid existence validation: pass"
	else
		echo "disk_cid existence validation: failed"
	fi
}

Check_authorization()
{
	if [ "$vm_cid" != "" -a "$disk_cid" != "" ]; then
		if [ `${curl_cmd}SoftLayer_Virtual_Guest/${vm_cid}/AttachedNetworkStorages/ISCSI.json | grep "\"id\":$disk_cid" |wc -l` == "1" ]; then
			echo "Authorization validation: pass"
		else
			echo "Authorization validation: failed"
		fi
	else
		echo "Authorization validation: failed"
	fi
}

Check_agent()
{
	vm_agentid=`sqlite3 $vm_pool_db "select agent_id from vms where private_ip=\"${director_ip}\";"`
	
	if [ "$vm_agentid" == "" ]; then
		echo "Agent ID validation: failed  Cannot find director ip $director_ip in DB" 
		exit 1
	else
		echo "Agent ID validation: pass"
		num=`grep -n disk_cid $bosh_deployment_yml | cut -d: -f 1`		
		line=":vm_agentid: $vm_agentid"
		sed -i "${num}a $line" $bosh_deployment_yml 
		sed -i 's/:vm_agentid:/  :vm_agentid:/g' $bosh_deployment_yml
	fi	
}

Change_stemcell_cid()
{
	sed -i "s/stemcell_cid: .*/stemcell_cid: '1234567'/g" $bosh_deployment_yml 
}

Pre-check
Get_ip
Check_vm
Check_disk
Check_authorization
Check_agent
Change_stemcell_cid
