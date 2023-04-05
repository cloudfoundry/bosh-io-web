# `bosh-io/web`

* [Publishing stemcells](docs/publishing-stemcells.md)

## Run Documentation Locally

### Prerequisites

* golang environment setup
* docker (if you want to build [docs-bosh](https://github.com/cloudfoundry/docs-bosh))

### Setup

```
./bin/build-docs.sh
./bin/git-init-clone.sh
./bin/run-local.sh
```

Open [localhost:3000](http://localhost:3000/)
