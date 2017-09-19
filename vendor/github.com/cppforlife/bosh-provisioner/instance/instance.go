package instance

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	bpdep "github.com/cppforlife/bosh-provisioner/deployment"
	bpinstupd "github.com/cppforlife/bosh-provisioner/instance/updater"
)

const instanceLogTag = "Instance"

type Instance struct {
	updater bpinstupd.Updater

	job         bpdep.Job
	depInstance bpdep.Instance

	logger boshlog.Logger
}

func NewInstance(
	updater bpinstupd.Updater,
	job bpdep.Job,
	depInstance bpdep.Instance,
	logger boshlog.Logger,
) Instance {
	return Instance{
		updater:     updater,
		job:         job,
		depInstance: depInstance,
		logger:      logger,
	}
}

func (i Instance) Deprovision() error {
	i.logger.Debug(instanceLogTag, "Tearing down instance")

	err := i.updater.TearDown()
	if err != nil {
		return bosherr.WrapErrorf(err, "Tearing down instance %d", i.depInstance.Index)
	}

	return nil
}
