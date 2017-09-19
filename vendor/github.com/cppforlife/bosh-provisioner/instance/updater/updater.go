package updater

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	bpeventlog "github.com/cppforlife/bosh-provisioner/eventlog"
	bpapplier "github.com/cppforlife/bosh-provisioner/instance/updater/applier"
)

const updaterLogTag = "Updater"

type Updater struct {
	instanceDesc string

	drainer     Drainer
	stopper     Stopper
	applier     bpapplier.Applier
	starter     Starter
	waiter      Waiter
	postStarter PostStarter

	eventLog bpeventlog.Log
	logger   boshlog.Logger
}

func NewUpdater(
	instanceDesc string,
	drainer Drainer,
	stopper Stopper,
	applier bpapplier.Applier,
	starter Starter,
	waiter Waiter,
	postStarter PostStarter,
	eventLog bpeventlog.Log,
	logger boshlog.Logger,
) Updater {
	return Updater{
		instanceDesc: instanceDesc,

		drainer:     drainer,
		stopper:     stopper,
		applier:     applier,
		starter:     starter,
		waiter:      waiter,
		postStarter: postStarter,

		eventLog: eventLog,
		logger:   logger,
	}
}

func (u Updater) SetUp() error {
	stage := u.eventLog.BeginStage(fmt.Sprintf("Setting up instance %s", u.instanceDesc), 3)

	task := stage.BeginTask("Applying")

	err := task.End(u.applier.Apply())
	if err != nil {
		return bosherr.WrapError(err, "Applying")
	}

	task = stage.BeginTask("Starting")

	err = task.End(u.starter.Start())
	if err != nil {
		return bosherr.WrapError(err, "Starting")
	}

	task = stage.BeginTask("Waiting")

	err = task.End(u.waiter.Wait())
	if err != nil {
		return bosherr.WrapError(err, "Waiting")
	}

	task = stage.BeginTask("Post-Start")

	err = task.End(u.postStarter.PostStart())
	if err != nil {
		return bosherr.WrapError(err, "Post-Starting")
	}

	return nil
}

func (u Updater) TearDown() error {
	stage := u.eventLog.BeginStage(fmt.Sprintf("Tearing down instance %s", u.instanceDesc), 2)

	task := stage.BeginTask("Draining")

	err := task.End(u.drainer.Drain())
	if err != nil {
		return bosherr.WrapError(err, "Draining")
	}

	task = stage.BeginTask("Stopping")

	err = task.End(u.stopper.Stop())
	if err != nil {
		return bosherr.WrapError(err, "Stopping")
	}

	return nil
}
