package noteimporter

import (
	"time"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

type FactoryOptions struct {
	Enabled bool
	Period  time.Duration

	GithubPersonalAccessToken string
}

type FactoryRepos interface {
	ReleasesRepo() bhrelsrepo.ReleasesRepository
}

type Factory struct {
	Importer NoteImporter
}

func NewFactory(
	options FactoryOptions,
	repos FactoryRepos,
	logger boshlog.Logger,
) (Factory, error) {
	if !options.Enabled {
		return Factory{Importer: NewNoopNoteImporter(logger)}, nil
	}

	periodicGithubNoteImporter := NewPeriodicGithubNoteImporter(
		options.Period,
		options.GithubPersonalAccessToken,
		make(chan struct{}),
		repos.ReleasesRepo(),
		logger,
	)

	return Factory{Importer: periodicGithubNoteImporter}, nil
}
