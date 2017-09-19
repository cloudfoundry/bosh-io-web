package updater

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	bpagclient "github.com/cppforlife/bosh-provisioner/agent/client"
)

const starterLogTag = "Starter"

type Starter struct {
	agentClient bpagclient.Client
	logger      boshlog.Logger
}

func NewStarter(
	agentClient bpagclient.Client,
	logger boshlog.Logger,
) Starter {
	return Starter{
		agentClient: agentClient,
		logger:      logger,
	}
}

func (s Starter) Start() error {
	s.logger.Debug(starterLogTag, "Running pre-start")

	if err := s.agentClient.PreStart(); err != nil {
		return bosherr.WrapError(err, "Pre-Starting")
	}

	s.logger.Debug(starterLogTag, "Starting instance")

	_, err := s.agentClient.Start()
	if err != nil {
		return bosherr.WrapError(err, "Starting")
	}

	return nil
}
