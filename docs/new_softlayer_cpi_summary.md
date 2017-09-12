# New softlayer CPI manifest specification

## manifest change a lot

This topic describes properties of deployment by new softlayer cpi by compare with the old SoftLayer CPI.

### AZs

AZs support the Director to eliminate and/or simplify manual configuration for balancing VMs across AZs and IP address management. Once AZs are defined, deployment jobs can be placed into one or more AZs.

AZs schema:  
  * azs [Array, required]: List of AZs.
  * name [String, required]: Name of an AZ within the Director.
  * cloud_properties [Hash, optional]: Describes any IaaS-specific properties needed to associated with AZ; for most IaaSes, some data here is actually required. See CPI Specific cloud_properties below. Example: availability_zone. Default is {} (empty Hash).

_Note that IaaS specific cloud properties related to AZs should now be only placed under azs. Make sure to remove them from resource_pools/vm_types’ cloud properties._

sample manifest for softlayer cpi:  
```yaml
azs:
- name: z1
  cloud_properties:
    Datacenter: { Name: lon02 }
```

sample manifest for new softlayer cpi:  
```yaml
azs:
- name: z1
  cloud_properties:
    datacenter: { name: lon02 }
```

### Networks

There are three different network types: `manual`, `dynamic`, and `vip`. And softlayer-cpi does not support `vip` type at present.

Manual networks schema:  
  * name [String, required]: Name used to reference this network configuration
  * type [String, required]: Value should be manual
  * subnets [Array, required]: Lists subnets in this network
    * range [String, required]: Subnet IP range that includes all IPs from this subnet
    * gateway [String, required]: Subnet gateway IP
    * dns [Array, optional]: DNS IP addresses for this subnet
    * reserved [Array, optional]: Array of reserved IPs and/or IP ranges. BOSH does not assign IPs from this range to any VM
    * static [Array, optional]: Array of static IPs and/or IP ranges. BOSH assigns IPs from this range to jobs requesting static IPs. Only IPs specified here can be used for static IP reservations.
    * az [String, optional]: AZ associated with this subnet (should only be used when using first class AZs). Example: z1. Available in v241+.
    * azs [Array, optional]: List of AZs associated with this subnet (should only be used when using first class AZs). Example: [z1, z2]. Available in v241+.
    * cloud_properties [Hash, optional]: Describes any IaaS-specific properties for the subnet. Default is {} (empty Hash).

sample manifest of static network for softlayer cpi:  
```yaml
networks:
- name: default
  type: manual
  subnets:
  - range: 10.112.166.128/26
    gateway: 10.112.166.129
    dns: [8.8.8.8, 10.0.80.11, 10.0.80.12]
    static: [10.112.166.131]
    az: z1
    cloud_properties:
      PrimaryNetworkComponent:
        NetworkVlan:
          Id: 524956
      PrimaryBackendNetworkComponent:
        NetworkVlan:
          Id: 524954
```

sample manifest of static network for new softlayer cpi:  
```yaml
networks:
- name: default
  type: manual
  subnets:
  - range: 10.112.166.128/26
    gateway: 10.112.166.129
    dns: [8.8.8.8, 10.0.80.11, 10.0.80.12]
    static: [10.112.166.131]
    az: z1
    cloud_properties:
      vlanIds: [524956, 524954]
```

Dynamic networks schema:  
  * name [String, required]: Name used to reference this network configuration
  * type [String, required]: Value should be dynamic
  * dns [Array, optional]: DNS IP addresses for this network
  * cloud_properties [Hash, optional]: Describes any IaaS-specific properties for the network. Default is {} (empty Hash).


sample manifest of dynamic network for softlayer cpi:  
```yaml
networks
- name: dynamic
  type: dynamic
  dns: [8.8.8.8, 10.0.80.11, 10.0.80.12]
  cloud_properties:
    PrimaryNetworkComponent:
      NetworkVlan:
        Id: 524956
    PrimaryBackendNetworkComponent:
      NetworkVlan:
        Id: 524954
```

sample manifest of dynamic network for new softlayer cpi:  
```yaml
- name: dynamic
  type: dynamic
  cloud_properties:
    vlanIds: [((sl_vlan_public_id)), ((sl_vlan_private_id))]
   dns: [8.8.8.8, 10.0.80.11, 10.0.80.12]
```

### VM Types/Resource Pools

VM type is a named Virtual Machine size configuration in the cloud config.
Resource pool is collections of VMs created from the same stemcell, with the same configuration, in a deployment.

Vm_types schema:  
  * vm_types [Array, required]: Specifies the VM types available to deployments. At least one should be specified.
    * name [String, required]: A unique name used to identify and reference the VM type
    * cloud_properties [Hash, optional]: Describes any IaaS-specific properties needed to create VMs; for most IaaSes, some data here is actually required. See CPI Specific cloud_properties below. Example: instance_type: m3.medium. Default is {} (empty Hash).

sample manifest of softlayer cpi: 
```yaml
vm_types:
- name: default
  cloud_properties:
    bosh_ip: ((internal_ip))
    startCpus:  4
    maxMemory:  8192
    ephemeralDiskSize: 100
    hourlyBillingFlag: true
    vmNamePrefix: manifest-sample
    domain: sofltayer.com
    datacenter:
      name: lon02
```

sample manifest for new softlayer cpi:  
```yaml
- name: default
  cloud_properties:
    startCpus:  4
    maxMemory:  8192
    maxNetworkSpeed: 100
    ephemeralDiskSize: 100
    hourlyBillingFlag: true
    vmNamePrefix: manifest-sample
    domain: sofltayer.com
    datacenter: lon02
```

### Disk types/Disk Pools

Disk Type (previously known as Disk Pool) is a named disk configuration specified in the cloud config.

* disk_types [Array, required]: Specifies the disk types available to deployments. At least one should be specified.
    * name [String, required]: A unique name used to identify and reference the disk type
    * disk_size [Integer, required]: Specifies the disk size. disk_size must be a positive integer. BOSH creates a persistent disk of that size in megabytes and attaches it to each job instance VM.
    * cloud_properties [Hash, optional]: Describes any IaaS-specific properties needed to create disks. Examples: type, iops. Default is {} (empty Hash).

```yaml
disk_types:
- name: disks
  disk_size: 100_000
  cloud_properties:
    Iops: 3000
    UseHourlyPricing: true
```

```yaml
disk_types:
- disk_size: 100000
  name: disks
  cloud_properties:
    iops: 3000
    snapshotSpace: 20

```
    
### more sample

director manifest comparison
```yaml
cloud_provider:
  cert:
    ca: ...
    certificate: ...
    private_key: ...
  mbus: https://mbus:lvwczko3gux720lijapg@10.112.166.130:6868
  properties:
    agent:
      mbus: https://mbus:lvwczko3gux720lijapg@0.0.0.0:6868
    blobstore:
      path: /var/vcap/micro_bosh/data/cache
      provider: local
    ntp:
    - time1.google.com
    - time2.google.com
    - time3.google.com
    - time4.google.com
    softlayer:
      apiKey: 7eab8fbfcdda3249e780dce0b10c7e4794e5ccd0fc9af7221b9fa9b40924ba8a
      username: cuixuex@cn.ibm.com
  template:
    name: softlayer_cpi
    release: bosh-softlayer-cpi
disk_pools:
- disk_size: 200000
  name: disks
instance_groups: ...
name: bosh
networks:
- name: default
  subnets:
  - cloud_properties:
      PrimaryNetworkComponent:
        NetworkVlan:
          Id: 1292653
    dns:
    - 8.8.8.8
    gateway: 10.112.166.129
    range: 10.112.166.128/26
    static:
    - 10.112.166.130
  type: manual
- cloud_properties:
    PrimaryBackendNetworkComponent:
      NetworkVlan:
        Id: 1292651
    PrimaryNetworkComponent:
      NetworkVlan:
        Id: 1292653
  dns:
  - 8.8.8.8
  - 10.0.80.11
  - 10.0.80.12
  name: dynamic
  type: dynamic
releases: ...
resource_pools:
- cloud_properties:
    datacenter:
      name: lon02
    deployedByBoshcli: true
    domain: softlayer.com
    ephemeralDiskSize: 100
    hourlyBillingFlag: true
    maxMemory: 8192
    networkComponents:
    - maxSpeed: 100
    startCpus: 4
    vmNamePrefix: director-sl-cpi-prefix
  env:
    bosh: ...
  name: vms
  network: dynamic
  stemcell: ...
variables: []
```

```yaml
cloud_provider:
  cert:
    ca: ...
    certificate: ...
    private_key: ...
  mbus: https://mbus:511kcx7pje314j5pk8dt@10.112.166.130:6868
  properties:
    agent:
      mbus: https://mbus:511kcx7pje314j5pk8dt@0.0.0.0:6868
    blobstore:
      path: /var/vcap/micro_bosh/data/cache
      provider: local
    ntp:
    - time1.google.com
    - time2.google.com
    - time3.google.com
    - time4.google.com
    softlayer:
      api_key: 7eab8fbfcdda3249e780dce0b10c7e4794e5ccd0fc9af7221b9fa9b40924ba8a
      username: cuixuex@cn.ibm.com
  ssh_tunnel:
    host: director-sl-cpi-shadow.softlayer.com
    port: 22
    private_key: ...
    user: root
  template:
    name: softlayer_cpi
    release: bosh-softlayer-cpi
disk_pools:
- cloud_properties:
    iops: 3000
    snapshotSpace: 20
  disk_size: 200000
  name: disks
instance_groups: ...
name: bosh
networks:
- name: default
  subnets:
  - cloud_properties:
      vlanIds:
      - 1292653
    dns:
    - 8.8.8.8
    gateway: 10.112.166.129
    range: 10.112.166.128/26
    static:
    - 10.112.166.130
  type: manual
- cloud_properties:
    vlanIds:
    - 1292651
    - 1292653
  dns:
  - 8.8.8.8
  - 10.0.80.11
  - 10.0.80.12
  name: dynamic
  type: dynamic
releases: ...
resource_pools:
- cloud_properties:
    datacenter: lon02
    deployedByBoshcli: true
    domain: softlayer.com
    ephemeralDiskSize: 100
    hourlyBillingFlag: true
    maxMemory: 8192
    maxNetworkSpeed: 100
    startCpus: 4
    vmNamePrefix: director-sl-cpi-prefix
  env:
    bosh: ...
  name: vms
  network: default
  stemcell: ...
variables: []
```


## diff
softlayer-go
implement actions
userData
registry
must settings
