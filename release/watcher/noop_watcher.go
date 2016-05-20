package watcher

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type NoopWatcher struct {
	logTag string
	logger boshlog.Logger
}

func NewNoopWatcher(logger boshlog.Logger) NoopWatcher {
	return NoopWatcher{
		logTag: "NoopWatcher",
		logger: logger,
	}
}

func (w NoopWatcher) Watch() error {
	w.logger.Info(w.logTag, "Starting noop watching")

	select {}
}
