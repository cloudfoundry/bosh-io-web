package updater

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	bpagclient "github.com/cppforlife/bosh-provisioner/agent/client"
)

const postStarterLogTag = "PostStarter"

type PostStarter struct {
	agentClient bpagclient.Client
	logger      boshlog.Logger
}

func NewPostStarter(
	agentClient bpagclient.Client,
	logger boshlog.Logger,
) PostStarter {
	return PostStarter{
		agentClient: agentClient,
		logger:      logger,
	}
}

// PostStart runs after an instance reaches running state.
func (w PostStarter) PostStart() error {
	w.logger.Debug(postStarterLogTag, "Running post-start")

	if err := w.agentClient.PostStart(); err != nil {
		return bosherr.WrapError(err, "Post-Starting")
	}

	return nil
}
