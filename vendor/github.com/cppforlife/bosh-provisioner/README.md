## Stand-alone BOSH provisioner

Stand-alone BOSH provisioner sets up and configures single VM
to be like any other BOSH managed VM. Besides installing
`bosh-agent` and `monit` on the system, BOSH provisioner
can optionally issue BOSH Agent's `apply` command to compile and
start running jobs as described by the deployment manifest.


### Usage

1. `go get github.com/cppforlife/bosh-provisioner/main` to install `bosh-provisioner`

2. `bosh-provisioner -configPath=./config.json` to run provisioner. Example `config.json`:

```
{
  assets_dir: "./assets",
  repos_dir: "/opt/bosh-provisioner/repos",

  blobstore: {
    provider: "local",
    options: {
      blobstore_path: "/opt/bosh-provisioner/blobstore",
    },
  },

  vm_provisioner: {
    full_stemcell_compatibility: false,

    agent_provisioner: {
      infrastructure: "warden",
      platform:       "ubuntu",
      configuration:  {},

      mbus: "https://user:password@127.0.0.1:4321/agent",
    },
  },

  deployment_provisioner: {
    manifest_path: "/opt/bosh-provisioner/manifest.yml",
  },
}
```

(Note: `assets_dir` includes pre-compiled assets for a default Ubuntu system.)
