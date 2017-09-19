package vagrant

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bpeventlog "github.com/cppforlife/bosh-provisioner/eventlog"
	bpvm "github.com/cppforlife/bosh-provisioner/vm"
)

type VMProvisionerFactory struct {
	fs     boshsys.FileSystem
	runner boshsys.CmdRunner

	assetsDir string
	mbus      string

	blobstoreConfig     map[string]interface{}
	vmProvisionerConfig bpvm.ProvisionerConfig

	eventLog bpeventlog.Log
	logger   boshlog.Logger
}

func NewVMProvisionerFactory(
	fs boshsys.FileSystem,
	runner boshsys.CmdRunner,
	assetsDir string,
	blobstoreConfig map[string]interface{},
	vmProvisionerConfig bpvm.ProvisionerConfig,
	eventLog bpeventlog.Log,
	logger boshlog.Logger,
) VMProvisionerFactory {
	return VMProvisionerFactory{
		fs:     fs,
		runner: runner,

		assetsDir:           assetsDir,
		blobstoreConfig:     blobstoreConfig,
		vmProvisionerConfig: vmProvisionerConfig,

		eventLog: eventLog,
		logger:   logger,
	}
}

func (f VMProvisionerFactory) NewVMProvisioner() *VMProvisioner {
	cmds := NewSimpleCmds(f.runner, f.logger)

	depsProvisionerFactory := NewDepsProvisionerFactory(
		f.vmProvisionerConfig.FullStemcellCompatibility,
		f.vmProvisionerConfig.AgentProvisioner.Platform,
		f.runner,
		f.eventLog,
		f.logger,
	)

	depsProvisioner := depsProvisionerFactory.NewDepsProvisioner()

	vcapUserProvisioner := NewVCAPUserProvisioner(
		f.fs,
		f.runner,
		f.eventLog,
		f.logger,
	)

	assetManager := NewAssetManager(f.assetsDir, f.fs, f.runner, f.logger)

	runitProvisioner := NewRunitProvisioner(
		f.fs,
		cmds,
		depsProvisioner,
		f.runner,
		assetManager,
		f.logger,
	)

	monitProvisioner := NewMonitProvisioner(
		cmds,
		assetManager,
		runitProvisioner,
		f.logger,
	)

	agentProvisioner := NewAgentProvisioner(
		f.fs,
		cmds,
		assetManager,
		runitProvisioner,
		monitProvisioner,
		f.blobstoreConfig,
		f.vmProvisionerConfig.AgentProvisioner,
		f.eventLog,
		f.logger,
	)

	vmProvisioner := NewVMProvisioner(
		vcapUserProvisioner,
		depsProvisioner,
		agentProvisioner,
		f.logger,
	)

	return vmProvisioner
}
