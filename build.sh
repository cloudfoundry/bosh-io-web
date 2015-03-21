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
# Use cppforlife/s3cli which now uses s3gof3r to take advantage of parallel uploads.
# (Had a problem uploading 5gb files with pivotal-golang/s3cli).
curl -L -o $tar https://github.com/cppforlife/s3cli/tarball/403686805a37d59babe5719568d78a16dbf8d8c4
mkdir $src_dir
tar -xzf $tar -C $src_dir/ --strip-components 1
( set -e; cd $src_dir; bin/build )
mv $src_dir/out/s3 ./bosh-blobstore-s3
rm -rf $tar $src_dir

echo "Building docs"
./build-docs.sh
