package vagrant

import (
	"fmt"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bpeventlog "github.com/cppforlife/bosh-provisioner/eventlog"
)

type DepsProvisionerFactory struct {
	fullStemcellCompatibility bool
	platform                  string

	runner   boshsys.CmdRunner
	eventLog bpeventlog.Log
	logger   boshlog.Logger
}

func NewDepsProvisionerFactory(
	fullStemcellCompatibility bool,
	platform string,
	runner boshsys.CmdRunner,
	eventLog bpeventlog.Log,
	logger boshlog.Logger,
) DepsProvisionerFactory {
	return DepsProvisionerFactory{
		fullStemcellCompatibility: fullStemcellCompatibility,
		platform:                  platform,

		runner:   runner,
		eventLog: eventLog,
		logger:   logger,
	}
}

func (f DepsProvisionerFactory) NewDepsProvisioner() DepsProvisioner {
	switch f.platform {
	case "ubuntu":
		return NewAptDepsProvisioner(
			f.fullStemcellCompatibility,
			f.runner,
			f.eventLog,
			f.logger,
		)
	case "centos":
		return NewYumDepsProvisioner(
			f.fullStemcellCompatibility,
			f.runner,
			f.eventLog,
			f.logger,
		)
	default:
		panic(fmt.Sprintf("Unknown dependency provisioner for platform '%s'", f.platform))
	}
}
