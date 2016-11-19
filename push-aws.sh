#!/bin/bash

set -e

ip=$1

echo "Pull down creds from lastpass to ./prod-conf"
rm -rf ./prod-conf && mkdir ./prod-conf
lpass show --notes "bosh-io-worker-json-conf" > ./prod-conf/worker.json
lpass show --notes "bosh-io-rds-combined-ca-bundle-pem" > ./prod-conf/rds-combined-ca-bundle.pem

config=prod-conf/worker.json

if [ ! -f $config ]; then
  echo "Missing $config config"
  exit 1
fi

./build.sh

echo "Pushing to $ip"
kill_marker=BOSH_HUB_KILL
src_path=/mnt/bosh-hub
log_path=/mnt/logs
tmp_path=/mnt/tmp

ssh -t vcap@$ip "
  set -e -x

  # Kill all processes with $kill_marker env var and without NO_$kill_marker
  sudo ps auxwwwe | grep $kill_marker | grep -v NO_$kill_marker | awk '{print "'$2'"}' | xargs kill -9 || true

  # Remove all dirs
  sudo rm -rf $src_path $log_path $tmp_path

  # Make sure dirs exist
  sudo mkdir -p $src_path $log_path $tmp_path

  # Make sure vcap has access to them
  sudo chown vcap:vcap $src_path $log_path $tmp_path
"

echo "Copying all needed files"
scp_files=(bosh-hub bosh-blobstore-s3 run.sh $config prod-conf/rds-combined-ca-bundle.pem Gemfile Gemfile.lock)
scp ${scp_files[*]} vcap@$ip:$src_path

rm -rf ./prod-conf

echo "Installing gems from Gemfile"
ssh -t vcap@$ip "cd $src_path && bundle config build.nokogiri --use-system-libraries && bundle install"

echo "Running new version of bosh-hub with BOSH_HUB_KILL env var"
nohup_cmd="nohup ./run.sh ./worker.json assets-id private-token-dont-matter >$log_path/stdout.log 2>$log_path/stderr.log"
ssh vcap@$ip "cd $src_path && $kill_marker=1 HOME=$tmp_path TMPDIR=$tmp_path exec $nohup_cmd &"

echo "Check if bosh-hub is running after 10 secs"
sleep 10
ssh vcap@$ip "pgrep bosh-hub"
