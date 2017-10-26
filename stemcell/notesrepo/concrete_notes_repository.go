package notesrepo

import (
	"encoding/json"
	"path/filepath"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type NotesRepository interface {
	Find(version string) (NoteRec, bool, error)
}

type NoteRec struct {
	Content string
}

type CNRepository struct {
	stemcellsLegacyIndexDir string
	stemcellsIndexDirs      []string

	fs     boshsys.FileSystem
	logger boshlog.Logger
}

func NewConcreteNotesRepository(
	stemcellsLegacyIndexDir string,
	stemcellsIndexDirs []string,
	fs boshsys.FileSystem,
	logger boshlog.Logger,
) CNRepository {
	return CNRepository{
		stemcellsLegacyIndexDir: stemcellsLegacyIndexDir,
		stemcellsIndexDirs:      stemcellsIndexDirs,
		fs:                      fs,
		logger:                  logger,
	}
}

func (r CNRepository) Find(version string) (NoteRec, bool, error) {
	var noteRec NoteRec

	contents, err := r.fs.ReadFile(filepath.Join(r.stemcellsLegacyIndexDir, version, "notes.v1.yml"))
	if err != nil {
		return noteRec, false, nil
	}

	err = json.Unmarshal(contents, &noteRec)
	if err != nil {
		return noteRec, false, bosherr.WrapError(err, "Unmarshaling stemcell notes")
	}

	return noteRec, true, nil
}
