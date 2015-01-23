### Setting up AWS instance to run worker

1. Run `terraform apply -var 'aws_access_key=XXX' -var 'aws_secret_key=XXX'`

2. Install `bundler` and `bosh_cli` gems

```
sudo apt-get -y install ruby ruby-dev
sudo rm -rf /var/vcap/bosh/bin
sudo gem install bundler bosh_cli --no-ri --no-rdoc
bosh -v
```

3. Attach a volume to the VM at `/dev/sdf`

4. Mount volume at `/mnt`

```
lsblk
sudo file -s /dev/xvdf
sudo mkfs -t ext4 /dev/xvdf
sudo mkdir /mnt
mount /dev/xvdf /mnt
sudo mount /dev/xvdf /mnt
```
