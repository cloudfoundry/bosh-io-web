package main

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bpeventlog "github.com/cppforlife/bosh-provisioner/eventlog"
	bpprov "github.com/cppforlife/bosh-provisioner/provisioner"
	bpvm "github.com/cppforlife/bosh-provisioner/vm"
)

var (
	DefaultConfig = Config{
		EventLog: bpeventlog.Config{
			DeviceType: bpeventlog.ConfigDeviceTypeJSON,
		},

		VMProvisioner: bpvm.ProvisionerConfig{
			AgentProvisioner: bpvm.AgentProvisionerConfig{
				Platform: "ubuntu",
				Mbus:     "https://user:password@127.0.0.1:4321/agent",
			},
		},
	}

	DefaultAgentConfiguration = map[string]interface{}{
		"Infrastructure": map[string]interface{}{
			"Settings": map[string]interface{}{
				"UseRegistry": true,
				"Sources": []map[string]interface{}{
					{
						"SettingsPath": "warden-cpi-agent-env.json",
						"Type":         "File",
					},
				},
			},
		},
		"Platform": map[string]interface{}{
			"Linux": map[string]interface{}{
				"UseDefaultTmpDir": true,
				"SkipDiskSetup":    true,
			},
		},
	}
)

type Config struct {
	// Assets dir is used as a temporary location to transfer files from host to guest.
	// It will not be created since assets already must be present.
	AssetsDir string `json:"assets_dir"`

	// Repos dir is mainly used to record what was placed in the blobstore.
	// It will be created if it does not exist.
	ReposDir string `json:"repos_dir"`

	// Tmp dir is used instead of the main tmp directory.
	// It will be created if it does not exist.
	TmpDir string `json:"tmp_dir"`

	EventLog bpeventlog.Config `json:"event_log"`

	Blobstore bpprov.BlobstoreConfig `json:"blobstore"`

	VMProvisioner bpvm.ProvisionerConfig `json:"vm_provisioner"`

	DeploymentProvisioner bpprov.DeploymentProvisionerConfig `json:"deployment_provisioner"`
}

func NewConfigFromPath(path string, fs boshsys.FileSystem) (Config, error) {
	var config Config

	bytes, err := fs.ReadFile(path)
	if err != nil {
		return config, bosherr.WrapErrorf(err, "Reading config %s", path)
	}

	config = DefaultConfig

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, bosherr.WrapError(err, "Unmarshalling config")
	}

	if config.VMProvisioner.AgentProvisioner.Configuration == nil {
		config.VMProvisioner.AgentProvisioner.Configuration = DefaultAgentConfiguration
	}

	err = config.validate()
	if err != nil {
		return config, bosherr.WrapError(err, "Validating config")
	}

	return config, nil
}

func (c Config) validate() error {
	if c.AssetsDir == "" {
		return bosherr.Error("Must provide non-empty assets_dir")
	}

	if c.ReposDir == "" {
		return bosherr.Error("Must provide non-empty repos_dir")
	}

	err := c.EventLog.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating event_log configuration")
	}

	if c.Blobstore.Type != bpprov.BlobstoreConfigTypeLocal {
		return bosherr.Error("Blobstore type must be local")
	}

	err = c.Blobstore.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating blobstore configuration")
	}

	return nil
}
