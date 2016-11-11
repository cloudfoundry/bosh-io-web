package notesrepo

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	bpindex "github.com/cppforlife/bosh-provisioner/index"
)

type CNRepository struct {
	index  bpindex.Index
	logger boshlog.Logger
}

type noteRecKey struct {
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

func (r CNRepository) Find(version string) (NoteRec, bool, error) {
	var noteRec NoteRec

	key := noteRecKey{VersionRaw: version}

	err := r.index.Find(key, &noteRec)
	if err != nil {
		if err == bpindex.ErrNotFound {
			return noteRec, false, nil
		}

		return noteRec, false, bosherr.WrapError(err, "Finding notes")
	}

	return noteRec, true, nil
}

func (r CNRepository) Save(version string, noteRec NoteRec) error {
	key := noteRecKey{VersionRaw: version}

	err := r.index.Save(key, noteRec)
	if err != nil {
		return bosherr.WrapError(err, "Saving notes")
	}

	return nil
}
