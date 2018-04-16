## `bosh-hub`

* [Publishing stemcells](docs/publishing-stemcells.md)

#### Run Documentation Locally

##### Prerequisites

* golang environment setup
* docker (if rebuilding docs-bosh)

##### Setup

```
go get github.com/bosh-io/web
cd $GOPATH/src/github.com/bosh-io/web
./build.sh
./run-local.sh

```

Open [http://localhost:3000/docs](http://localhost:3000/docs)
