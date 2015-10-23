
# Deploying a Concourse Pipeline on SoftLayer
As the title suggests, this is a guide to deploying concourse on a SoftLayer VM. While it isn't the shortest guide, it is relatively straightforward, so hopefully it can be of some help until an easier solution is created.  

If you have an existing VM you'd like to use, skip the following section and go directly to the one titled [Linking Vagrant](#link).  
Whichever way you decide to create a VM, I'd recommend you create a virtual machine with a primary disk larger than the base 25GB. The storage requirements for some of our tasks are surprisingly large, given the amount of dependencies required, and adding on storage after the fact (i.e. on a second physical disk using LVM) can be frustrating. Similarly, it may be prudent to add additional processing capability and memory to the VM.

<h2 id="createVM">Creating a SoftLayer Instance</h2>
The [Vagrant SoftLayer Provider](https://github.com/audiolize/vagrant-softlayer) helps you create a SoftLayer VM using Vagrant. Follow the README on their github page to get started. The gist of it is:

* Install the vagrant-softlayer plugin:
  ```
  $ vagrant plugin install vagrant-softlayer
  ```

* Initialize vagrant:
  ```
  $ vagrant init softlayer.box
  ```
  - You can create your own box file by following the README [here](https://github.com/audiolize/vagrant-softlayer/tree/master/example_box).  
    Learn more about the box file format [here](http://docs.vagrantup.com/v2/boxes/format.html).

* Edit the `Vagrantfile` to reflect your specific provider needs:
  ```ruby
  Vagrant.configure(2) do |config|
    # ... other stuff
  
    config.vm.provider :softlayer do |sl, override|
      sl.domain = "whatever.com"
      sl.ssh_key = "ssh_key_name"
      sl.hostname = "something"
      sl.datacenter = "ams01"
      override.ssh.username = "root"
      override.ssh.private_key_path = "/your/private_key_path"
    end
  end
  ```

* Create and configure guest machine using the `softlayer` provider in your `Vagrantfile`:
  ```
  $ vagrant up --provider=softlayer
  ```

<h2 id="link">Linking Vagrant</h2>

The [Vagrant ManagedServers Provider](https://github.com/tknerr/vagrant-managed-servers) allows you to "link" vagrant with an existing managed server. We're going to use it to link to our existing SoftLayer instance.  

* Why don't we just provision using the existing `softlayer` provider? Well, we get the following error:
  ```
  No host IP was given to the Vagrant core NFS helper. This is 
  an internal error that should be reported as a bug.
  ```
  
  So maybe that's something worth exploring further.  
  
Anyway, to continue:

* Install using the standard Vagrant plugin installation method:
  ```
  $ vagrant plugin install vagrant-managed-servers
  ```

* The following steps differ depending on if you created a SoftLayer instance using the Vagrant Softlayer Provider:
  - If you skipped the previous section to [create a SoftLayer instance]($createVM), go ahead and create a `Vagrantfile` in your current directory:
    ```
    $ vagrant init tknerr/managed-server-dummy
    ```

  - If you didn't skip the previous section, you'll have to delete the `.vagrant` directory in order for you to link the existing VM to the managed provider. You'll also have to edit the `Vagrantfile` to reflect a different box:
    ```ruby
    Vagrant.configure("2") do |config|
      config.vm.box = "tknerr/managed-server-dummy"
      # ... other stuff
    end
    ```

* In the `Vagrantfile` you can now use the `managed` provider and specify the managed server's hostname and credentials:
  ```ruby
  Vagrant.configure("2") do |config|
    # ... other stuff
    
    config.vm.provider :managed do |managed, override|
      managed.server = "something.whatever.com"
      override.ssh.username = "root"
      override.ssh.private_key_path = "/your/private_key_path"
    end
  end
```

* Link the Vagrant VM with your managed server by running:
  ```
  $ vagrant up --provider=managed
  ```
  - Unfortunately, we can't just run this command to begin with because Vagrant doesn't currently support multiple providers:
    ```
    An active machine was found with a different provider. Vagrant
    currently allows each machine to be brought up with only a single
    provider at a time. A future version will remove this limitation.
    Until then, please destroy the existing machine to up with a new
    provider.

    Machine name: default
    Active provider: softlayer
    Requested provider: managed
    ```

## Provisioning the VM
The [BOSH provisioner for Vagrant](https://github.com/cppforlife/vagrant-bosh) lets you provision guest VMs by specifying a regular BOSH deployment manifest.

* Install as before using `vagrant plugin install`:
  ```shell
  $ vagrant plugin install vagrant-bosh
  
  ```

* Now, create a manifest file for use with provisioning. In this example, we create a file called `concourse_manifest.yml`:
  ```yml
  ---
  name: concourse

  releases: # you may have to update these attributes as new releases become available
    - name: concourse
      url: https://github.com/concourse/concourse/releases/download/v0.65.1/concourse-0.65.1.tgz
      version: 0.65.1
    - name: garden-linux
      url: https://github.com/concourse/concourse/releases/download/v0.65.1/garden-linux-0.307.0.tgz
      version: 0.307.0

  networks:
    - name: concourse
      type: dynamic

  jobs:
    - name: concourse
      instances: 1
      networks: [{name: concourse}]
      templates:
        - {release: concourse, name: atc}
        - {release: concourse, name: tsa}
        - {release: concourse, name: groundcrew}
        - {release: concourse, name: postgresql}
        - {release: garden-linux, name: garden}
      properties:
        atc:
          publicly_viewable: true # anyone can view the main pipeline page

          postgresql:
            address: 127.0.0.1:5432
            role: &atc-role
              name: atc
              password: dummy-password

          basic_auth_username: SOME_USERNAME # set to lock builds
          basic_auth_password: SOME_PASSWORD

        postgresql:
          databases: [{name: atc}]
          roles:
            - *atc-role

        tsa:
          atc:
            address: 127.0.0.1:8080
            username: SOME_USERNAME # use same credentials as above
            password: SOME_PASSWORD

        groundcrew:
          tsa:
            host: 127.0.0.1

          garden:
            address: 127.0.0.1:7777

        garden:
          disk_quota_enabled: false # disk quotas are not enabled by default on SoftLayer VMs

          listen_network: tcp
          listen_address: 0.0.0.0:7777

          allow_host_access: true

  compilation:
    network: concourse

  update:
    canaries: 0
    canary_watch_time: 1000-60000
    update_watch_time: 1000-60000
    max_in_flight: 10
  ```
  - [Here](http://concourse.ci/deploying-and-upgrading-concourse.html#bosh-properties)'s a list of common values you might want to change for your particular implementation.
  
* We'll have to edit the `Vagrantfile` again to let the provisioner know where to find the file we just created:
  ```ruby
  Vagrant.configure(2) do |config|
    # ... other stuff
    
    config.vm.provision "bosh" do |c|
      # use cat or just inline full deployment manifest
      c.manifest = File.read(File.expand_path("../concourse_manifest.yml", __FILE__))
    end
  end
  ```

* Provision the VM, and you're done:
  ```
  $ vagrant provision
  ```

## Using Concourse
Concourse docs can be found at http://concourse.ci/  
View an example pipeline YAML in [concourse.yml](concourse.yml)

## Troubleshooting
* git resource can't be accessed
  * Make sure you've passed your private key.
  * If you have, it may be that the default SoftLayer DNS servers don't allow you to resolve outside domain names, such as `github.com`. If you find this to be the case (such as by `ping`ing), manually edit `etc/resolv.conf` in your SoftLayer instance (`vagrant ssh`) to include another DNS server, such as Google's `8.8.8.8`.

* no space left on the device (whether you checked this in the vm itself or saw an error on a concourse job)
  * Check the space left on your vm with `df -h`
  * If there is no space left on `/dev/xvda2`:
    * `monit stop garden`
    * `rm -rf /var/vcap/data/garden`
      * wait until process finishes
    * `monit start garden`

* general troubleshooting
  * Use `monit status` in the concourse vm to monitor the status of individual process
    * Use something like `watch monit status` if you want to automatically refresh it on an interval (such as if you're waiting for the processes to initialize) 
  * You can view logs in the concourse vm at `/var/vcap/data/sys/log`
    * If you get an error involving workers not being available, try the `groundcrew` or `postgres` logs.
  
