bosh-softlayer-tools [![Build Status](https://travis-ci.org/cloudfoundry-community/bosh-softlayer-tools.svg?branch=master)](https://travis-ci.org/cloudfoundry-community/bosh-softlayer-tools#)
========================

CLIs and other tools to ease in creating [BOSH](https://github.com/cloudfoundry/bosh) stemcells (light) for SoftLayer as well as manage pools of SoftLayer bare metal VM. All designed to work with the [bosh-softlayer-cpi](https://github.com/cloudfoundry/bosh-softlayer-cpi) project.

## Getting Started (*)
------------------

TBD

### Overview Presentations (*)
--------------------------

TBD

### Cloning and Building
------------------------

Clone this repo and build it. Using the following commands on a Linux or Mac OS X system:

```
$ mkdir -p bosh-softlayer-tools/src/cloudfoundry-community/bosh-softlayer-tools
$ export GOPATH=$(pwd)/bosh-softlayer-tools:$GOPATH
$ cd bosh-softlayer-tools/src/cloudfoundry-community/bosh-softlayer-tools
$ git clone https://cloudfoundry-community/bosh-softlayer-tools/bosh-softlayer-tools.git
$ cd bosh-softlayer-tools
$ ./bin/build
$ ./bin/test-unit
$ export SL_USERNAME=your-username@your-org.com
$ export SL_API_KEY=your-softlayer-api-key
$ ./bin/test-integration
```

NOTE: if you get any dependency errors, then use `go get path/to/dependency` to get it, e.g., `go get github.com/onsi/ginkgo` and `go get github.com/onsi/gomega`

The executables output should now be located in: `out/stemcells` and `out\bmp`.

### Running Tests
-----------------

The [SoftLayer](http://www.softlayer.com) (SL) stemcells, and associated tests and binary distribution depend on you having a real SL account. Get one for free for one month [here](http://www.softlayer.com/info/free-cloud). From your SL account you can get an API key. Using your account name and API key you will need to set two environment variables: `SL_USERNAME` and `SL_API_KEY`. You can do so as follows:

```
$ export SL_USERNAME=your-username@your-org.com
$ export SL_API_KEY=your-softlayer-api-key
```

You should run the tests to make sure all is well, do this with: `$ ./bin/test-unit` and `$ ./bin/test-integration` in your cloned repository. Please note that the `$ ./bin/test-integration` will spin up real SoftLayer virtual guests (VMs) and associated resources and will also delete them. This integration test may take up to 30 minutes (usually shorter)

The output should of `$ ./bin/test-unit` be similar to:

```
$ ./bin/build
$ ./bin/test-unit

 Cleaning build artifacts...

 Formatting packages...

 Integration Testing packages:
Will skip:
  ./integration
  [...]
Test Suite Passed

 Vetting packages for potential issues...

SWEET SUITE SUCCESS
```

## Developing
-------------

1. Check for existing stories on our [public Tracker](https://www.pivotaltracker.com/n/projects/1344876)
2. Select an unstarted story and work on code for it
3. If the story you want to work on is not there then open an issue and ask for a new story to be created
4. Run `go get golang.org/x/tools/cmd/vet`
5. Run `go get github.com/xxx ...` to install test dependencies (as you see errors)
6. Write a [Ginkgo](https://github.com/onsi/ginkgo) test.
7. Run `bin/test` and watch the test fail.
8. Make the test pass.
9. Submit a pull request.

## Contributing
---------------

* We gratefully acknowledge and thank the [current contributors](https://hithub.com/cloudfoundry-community/bosh-softlayer-tools/graphs/contributors)
* We welcome any and all contributions as Pull Requests (PR)
* We also welcome issues and bug report and new feature request. We will address as time permits
* Follow the steps above in Developing to get your system setup correctly
* Please make sure your PR is passing Travis before submitting
* Feel free to email me or the current collaborators if you have additional questions about contributions
* Before submitting your first PR, please read and follow steps in [CONTRIBUTING.md](CONTRIBUTING.md)

### Managing dependencies
-------------------------

* All dependencies managed via [Godep](https://github.com/tools/godep). See [Godeps/_workspace](https://github.com/cloudfoundry-community/bosh-softlayer-tools/tree/master/Godeps/_workspace) directory on master

#### Short `godep` Guide
* If you ever import a new package `foo/bar` (after you `go get foo/bar`, so that foo/bar is in `$GOPATH`), you can type `godep save ./...` to add it to the `Godeps` directory.
* To restore dependencies from the `Godeps` directory, simply use `godep restore`. `restore` is the opposite of `save`.
* If you ever remove a dependency or a link becomes deprecated, the easiest way is probably to remove your entire `Godeps` directory and run `godep save ./...` again, after making sure all your dependencies are in your `$GOPATH`. Don't manually edit `Godeps.json`!
* To update an existing dependency, you can use `godep update foo/bar` or `godep update foo/...` (where `...` is a wildcard)
* The godep project [readme](https://github.com/tools/godep/README.md) is a pretty good resource: [https://github.com/tools/godep](https://github.com/tools/godep)

* Since GO1.5, dependencies is managed via [Govendor](https://github.com/kardianos/govendor). See [vendor](https://github.com/cloudfoundry/bosh-softlayer-cpi/tree/master/vendor) directory.

### Current conventions
-----------------------

* Basic Go conventions
* Strict TDD for any code added or changed
* Go fakes when needing to mock objects

(*) these items are in the works, we will remove the * once the section is complete.

**NOTE**: this BOSH light stemcell project is used to support the [bosh-softlayer-cpi](https://github.com/cloudfoundry/bosh-softlayer-cpi) project. Like the CPI project, consider this code as prototype code.
