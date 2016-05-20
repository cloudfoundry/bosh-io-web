package importer

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type NoopImporter struct {
	logTag string
	logger boshlog.Logger
}

func NewNoopImporter(logger boshlog.Logger) NoopImporter {
	return NoopImporter{
		logTag: "NoopImporter",
		logger: logger,
	}
}

func (i NoopImporter) Import() error {
	i.logger.Info(i.logTag, "Starting noop importing bosh-init-bins")

	select {}
}
