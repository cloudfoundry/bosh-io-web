### Using `bosh-blobstore-s3`

Running `bosh-blobstore-s3` command:

```
bosh-blobstore-s3 -c ~/workspace/s3-config.json put ~/workspace/concourse/concourse-0.16.0.tgz some-tarball-guid
```

with `s3-config.json`:

```json
{
  "access_key_id":"AKIA...",
  "secret_access_key":"WG3L...",
  "bucket_name":"bosh-hub-..."
}
```
