package stemsrepo

import (
	"encoding/json"
	"path/filepath"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bhnotesrepo "github.com/cppforlife/bosh-hub/stemcell/notesrepo"
)

type S3StemcellsRepository struct {
	stemcellsIndexDir string
	notesRepo         bhnotesrepo.NotesRepository
	fs                boshsys.FileSystem
	logger            boshlog.Logger
}

type s3StemcellRec struct {
	Key  string
	ETag string
	SHA1 string

	Size         uint64
	LastModified string

	URL string
}

func NewS3StemcellsRepository(
	stemcellsIndexDir string,
	notesRepo bhnotesrepo.NotesRepository,
	fs boshsys.FileSystem,
	logger boshlog.Logger,
) S3StemcellsRepository {
	return S3StemcellsRepository{
		stemcellsIndexDir: stemcellsIndexDir,
		notesRepo:         notesRepo,
		fs:                fs,
		logger:            logger,
	}
}

func (r S3StemcellsRepository) FindAll(name string) ([]Stemcell, error) {
	var stems []Stemcell

	var s3StemcellRecs []s3StemcellRec

	contents, err := r.fs.ReadFile(filepath.Join(r.stemcellsIndexDir, "index.json"))
	if err != nil {
		return nil, bosherr.WrapError(err, "Finding all S3 stemcell recs")
	}

	err = json.Unmarshal(contents, &s3StemcellRecs)
	if err != nil {
		return nil, bosherr.WrapError(err, "Unmarshaling S3 stemcell recs")
	}

	filterByName := len(name) > 0

	for _, rec := range s3StemcellRecs {
		stemcell := NewS3Stemcell(
			rec.Key,
			rec.ETag,
			rec.SHA1,
			rec.Size,
			rec.LastModified,
			rec.URL,
		)

		if stemcell == nil || stemcell.IsDeprecated() {
			continue
		}

		stemcell.notesRepo = r.notesRepo

		if !filterByName || filterByName && stemcell.Name() == name {
			stems = append(stems, stemcell)
		}
	}

	return stems, nil
}
