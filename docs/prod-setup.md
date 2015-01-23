### Prod setup

Since CF apps are limited to 2gb of disk space not all releases can be built (e.g. cf-release is about 3.6gb). Once CF supports custom disk size configuration bosh-hub can be fully deployed to CF.

Currently app has to be deployed as two components:
- web on CF as `bosh-hub` CF app to `cfcommunity` organization, `boshorg` space
- worker on AWS

Both components use the same code, but they are configured differently via the JSON file. Since worker does not serve traffic it uses ActAsWorker=true and release and stemcell importers enabled. The web component serves web traffic and has importing disabled.

Both components have access to the Postgres DB running on AWS RDS. Since the worker is the only component that uploads generated releases, the web component does not have access to the S3 blobstore.

Before pushing run `./get-subs.sh`

#### Push web

Login to your cf account with cf CLI and run `./push-cf.sh`.

#### Push worker

Add SSH key to the keychain and run `./push-aws.sh IP`
