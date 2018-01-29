## `bosh-hub`

* [Tech overview](docs/tech-overview.md)
* [Prod setup](docs/prod-setup.md)
* [Worker setup](docs/worker-setup.md)

#### Misc

* [Recorded releases](docs/recorded-releases.md)
* [Using bosh-blobstore-s3](docs/using-bosh-blobstore-s3.md)

#### Run Documentation Locally

##### Prerequisites

* latest version of ruby 2.0.0 
* golang environment setup

##### Setup

```
go get github.com/bosh-io/web
cd $GOPATH/src/github.com/bosh-io/web
git clone https://github.com/cloudfoundry/docs-bosh.git
cd docs-bosh-io
bundle install 
#`bundle update <gem>` can help in case `bundle install` has issues with installing dependencies  
cd ..
./build-docs.sh
./run-local.sh

```

Open [http://localhost:3000/docs](http://localhost:3000/docs)
 

