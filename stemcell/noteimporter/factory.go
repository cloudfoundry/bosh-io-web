package noteimporter

import (
	"time"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	bhstemnotesrepo "github.com/cppforlife/bosh-hub/stemcell/notesrepo"
)

type FactoryOptions struct {
	Enabled bool
	Period  time.Duration

	GithubPersonalAccessToken string
}

type FactoryRepos interface {
	StemcellNotesRepo() bhstemnotesrepo.NotesRepository
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
		repos.StemcellNotesRepo(),
		logger,
	)

	return Factory{Importer: periodicGithubNoteImporter}, nil
}
