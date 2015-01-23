package watcher

import (
	"time"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	bhfetcher "github.com/cppforlife/bosh-hub/release/fetcher"
	bhimpsrepo "github.com/cppforlife/bosh-hub/release/importsrepo"
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhwatchersrepo "github.com/cppforlife/bosh-hub/release/watchersrepo"
)

type FactoryOptions struct {
	Enabled bool
	Period  time.Duration
}

type FactoryRepos interface {
	ReleasesRepo() bhrelsrepo.ReleasesRepository
	ImportsRepo() bhimpsrepo.ImportsRepository
	WatchersRepo() bhwatchersrepo.WatchersRepository
}

type Factory struct {
	Watcher Watcher
}

func NewFactory(
	options FactoryOptions,
	repos FactoryRepos,
	fetcher bhfetcher.Fetcher,
	logger boshlog.Logger,
) (Factory, error) {
	if !options.Enabled {
		return Factory{Watcher: NewNoopWatcher(logger)}, nil
	}

	periodicWatcher := NewPeriodicWatcher(
		options.Period,
		make(chan struct{}),
		repos.ReleasesRepo(),
		repos.WatchersRepo(),
		repos.ImportsRepo(),
		fetcher,
		logger,
	)

	return Factory{Watcher: periodicWatcher}, nil
}
