package instance

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	bpagclient "github.com/cppforlife/bosh-provisioner/agent/client"
	bpdep "github.com/cppforlife/bosh-provisioner/deployment"
	bpinstupd "github.com/cppforlife/bosh-provisioner/instance/updater"
)

const provisionerLogTag = "Provisioner"

type Provisioner struct {
	instanceUpdaterFactory bpinstupd.Factory
	logger                 boshlog.Logger
}

func NewProvisioner(
	instanceUpdaterFactory bpinstupd.Factory,
	logger boshlog.Logger,
) Provisioner {
	return Provisioner{
		instanceUpdaterFactory: instanceUpdaterFactory,
		logger:                 logger,
	}
}

func (p Provisioner) Provision(ac bpagclient.Client, job bpdep.Job, depInstance bpdep.Instance) (Instance, error) {
	p.logger.Debug(provisionerLogTag, "Updating instance")

	updater := p.instanceUpdaterFactory.NewUpdater(ac, job, depInstance)

	err := updater.SetUp()
	if err != nil {
		return Instance{}, bosherr.WrapErrorf(err, "Updating instance %d", depInstance.Index)
	}

	return NewInstance(updater, job, depInstance, p.logger), nil
}

func (p Provisioner) PreviouslyProvisioned(ac bpagclient.Client, job bpdep.Job, depInstance bpdep.Instance) Instance {
	p.logger.Debug(provisionerLogTag, "Finding previously provisioned instance")

	updater := p.instanceUpdaterFactory.NewUpdater(ac, job, depInstance)

	return NewInstance(updater, job, depInstance, p.logger)
}
