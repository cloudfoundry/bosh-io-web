package noteimporter

import (
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
)

type NoopNoteImporter struct {
	logTag string
	logger boshlog.Logger
}

func NewNoopNoteImporter(logger boshlog.Logger) NoopNoteImporter {
	return NoopNoteImporter{
		logTag: "stemcell.NoopNoteImporter",
		logger: logger,
	}
}

func (i NoopNoteImporter) Import() error {
	i.logger.Debug(i.logTag, "Starting noop note importing")

	select {}
}
