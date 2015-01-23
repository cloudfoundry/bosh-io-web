#!/bin/bash

set -e

echo "Configuring go"
export GOARCH=amd64
export GOOS=linux
export GOTOOLDIR=$(go env GOROOT)/pkg/linux_amd64

echo "Building bosh-hub"
go build -o bosh-hub github.com/cppforlife/bosh-hub/main

echo "Building bosh-blobstore-s3"
tar=s3cli.tar.gz
src_dir=s3cli-src
rm -rf $src_dir
curl -L -o $tar https://github.com/pivotal-golang/s3cli/tarball/2c4a7f0ceef411532bb051e7ca55a490a565cf60
mkdir $src_dir
tar -xzf $tar -C $src_dir/ --strip-components 1
( set -e; cd $src_dir; bin/build )
mv $src_dir/out/s3 ./bosh-blobstore-s3
rm -rf $tar $src_dir

echo "Building docs"
./build-docs.sh
