package notesrepo

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	bhrelver "github.com/cppforlife/bosh-hub/release/relver"
)

type NotesRepository interface {
	Find(source, version string) (NoteRec, bool, error)
}

type NoteRec struct {
	Content string
}

type CNRepository struct {
	relVerFactory bhrelver.Factory
	logger        boshlog.Logger
}

func NewConcreteNotesRepository(
	relVerFactory bhrelver.Factory,
	logger boshlog.Logger,
) CNRepository {
	return CNRepository{
		relVerFactory: relVerFactory,
		logger:        logger,
	}
}

func (r CNRepository) Find(source, version string) (NoteRec, bool, error) {
	var noteRec NoteRec

	relVer, err := r.relVerFactory.Find(source, version)
	if err != nil {
		return noteRec, false, err
	}

	contents, found, err := relVer.ReadOptinal("notes.v1.yml")
	if err != nil {
		return noteRec, false, err
	}

	if found {
		err = json.Unmarshal(contents, &noteRec)
		if err != nil {
			return noteRec, false, bosherr.WrapError(err, "Unmarshaling release notes")
		}

		return noteRec, true, nil
	}

	return noteRec, false, nil
}
