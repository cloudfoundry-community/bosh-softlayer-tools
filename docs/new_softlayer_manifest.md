# New softlayer manifest

## cf-deployment manifest

```
---
name: <%= properties.name || "bat" %>

releases:
  - name: bat
    version: <%= properties.release || "latest" %>

compilation:
  workers: 1
  network: default
  reuse_compilation_vms: true
  cloud_properties:
    startCpus: 4
    maxMemory: 8192
    maxNetworkSpeed: 100
    ephemeralDiskSize: 100
    datacenter: <%= properties.cloud_properties.data_center %>
    vmNamePrefix: bat-softlayer-worker
    domain: softlayer.com
    hourlyBillingFlag: true
    localDiskFlag: false

update:
  canaries: <%= properties.canaries || 1 %>
  canary_watch_time: 3000-90000
  update_watch_time: 3000-90000
  max_in_flight: <%= properties.max_in_flight || 1 %>
  serial: true

networks:
<% properties.networks.each do |network| %>
- name: <%= network.name %>
  type: dynamic
  dns: <%= properties.dns %>
  cloud_properties:
    vlanIds: <%= network.cloud_properties.vlanIds %>
<% end %>

resource_pools:
  - name: common
    network: default
    size: <%= properties.pool_size %>
    stemcell:
      name: <%= properties.stemcell.name %>
      version: '<%= properties.stemcell.version %>'
    cloud_properties:
      startCpus: 4
      maxMemory: 8192
      maxNetworkSpeed: 100
      ephemeralDiskSize: 100
      datacenter: <%= properties.cloud_properties.data_center %>
      vmNamePrefix: <%= properties.cloud_properties.vm_name_prefix %>
      domain: <%= properties.cloud_properties.domain %>
      hourlyBillingFlag: true
      localDiskFlag: false
    <% if properties.password %>
    env:
      bosh:
        password: <%= properties.password %>
        keep_root_password: true
    <% end %>

jobs:
  - name: <%= properties.job || "batlight" %>
    templates: <% (properties.templates || ["batlight"]).each do |template| %>
    - name: <%= template %>
    <% end %>
    instances: <%= properties.instances %>
    resource_pool: common
    <% if properties.persistent_disk %>
    persistent_disk: <%= properties.persistent_disk %>
    <% end %>
    networks:
      - name: default

properties:
  batlight:
    <% if properties.batlight.fail %>
    fail: <%= properties.batlight.fail %>
    <% end %>
    <% if properties.batlight.missing %>
    missing: <%= properties.batlight.missing %>
    <% end %>
    <% if properties.batlight.drain_type %>
    drain_type: <%= properties.batlight.drain_type %>
    <% end %>
```