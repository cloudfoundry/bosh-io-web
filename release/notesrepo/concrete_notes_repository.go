package notesrepo

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	bpindex "github.com/cppforlife/bosh-provisioner/index"
)

type CNRepository struct {
	index  bpindex.Index
	logger boshlog.Logger
}

type noteRecKey struct {
	Source     string
	VersionRaw string
}

func NewConcreteNotesRepository(
	index bpindex.Index,
	logger boshlog.Logger,
) CNRepository {
	return CNRepository{
		index:  index,
		logger: logger,
	}
}

func (r CNRepository) Find(source, version string) (NoteRec, bool, error) {
	var noteRec NoteRec

	key := noteRecKey{Source: source, VersionRaw: version}

	err := r.index.Find(key, &noteRec)
	if err != nil {
		if err == bpindex.ErrNotFound {
			return noteRec, false, nil
		}

		return noteRec, false, bosherr.WrapError(err, "Finding notes")
	}

	return noteRec, true, nil
}

func (r CNRepository) Save(source, version string, noteRec NoteRec) error {
	key := noteRecKey{Source: source, VersionRaw: version}

	err := r.index.Save(key, noteRec)
	if err != nil {
		return bosherr.WrapError(err, "Saving notes")
	}

	return nil
}
