package stemsrepo

import (
	"encoding/json"
	"encoding/xml"
	"path/filepath"
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bhnotesrepo "github.com/bosh-io/web/stemcell/notesrepo"

	"github.com/dpb587/metalink"
)

type S3StemcellsRepository struct {
	legacyStemcellsIndexDir string
	stemcellsIndexDirs      []string

	notesRepo bhnotesrepo.NotesRepository
	fs        boshsys.FileSystem
	logger    boshlog.Logger
}

type s3StemcellRec struct {
	Key    string
	ETag   string
	SHA1   string
	SHA256 string

	Size         uint64
	LastModified string

	URL string
}

func NewS3StemcellsRepository(
	legacyStemcellsIndexDir string,
	stemcellsIndexDirs []string,
	notesRepo bhnotesrepo.NotesRepository,
	fs boshsys.FileSystem,
	logger boshlog.Logger,
) S3StemcellsRepository {
	return S3StemcellsRepository{
		legacyStemcellsIndexDir: legacyStemcellsIndexDir,
		stemcellsIndexDirs:      stemcellsIndexDirs,
		notesRepo:               notesRepo,
		fs:                      fs,
		logger:                  logger,
	}
}

func (r S3StemcellsRepository) FindAll(name string) ([]Stemcell, error) {
	var stems []Stemcell

	s3StemcellRecs, err := r.loadLegacyIndex()
	if err != nil {
		return nil, bosherr.WrapError(err, "Loading stemcells from legacy index")
	}

	for _, dir := range r.stemcellsIndexDirs {
		recs, err := r.loadIndex(dir)
		if err != nil {
			return nil, bosherr.WrapErrorf(err, "Loading stemcells from index '%s'", dir)
		}

		s3StemcellRecs = append(s3StemcellRecs, recs...)
	}

	filterByName := len(name) > 0

	for _, rec := range s3StemcellRecs {
		err := rec.Validate()
		if err != nil {
			return nil, bosherr.WrapError(err, "Validating stemcell record")
		}

		stemcell := NewS3Stemcell(
			rec.Key,
			rec.ETag,
			rec.SHA1,
			rec.SHA256,
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

func (r S3StemcellsRepository) loadLegacyIndex() ([]s3StemcellRec, error) {
	recs := []s3StemcellRec{}

	contents, err := r.fs.ReadFile(filepath.Join(r.legacyStemcellsIndexDir, "index.json"))
	if err != nil {
		return nil, bosherr.WrapError(err, "Finding all S3 stemcell recs")
	}

	err = json.Unmarshal(contents, &recs)
	if err != nil {
		return nil, bosherr.WrapError(err, "Unmarshaling S3 stemcell recs")
	}

	return recs, nil
}

func (r S3StemcellsRepository) loadIndex(dir string) ([]s3StemcellRec, error) {
	recs := []s3StemcellRec{}

	foundPaths, err := r.fs.Glob(filepath.Join(dir, "published", "*", "*", "*.meta4"))
	if err != nil {
		return nil, bosherr.WrapError(err, "Globbing stemcells")
	}

	for _, path := range foundPaths {
		contents, err := r.fs.ReadFile(path)
		if err != nil {
			return nil, bosherr.WrapError(err, "Reading release version")
		}

		var meta4 metalink.Metalink

		err = xml.Unmarshal(contents, &meta4)
		if err != nil {
			return nil, bosherr.WrapError(err, "Unmarshaling meta4")
		}

		for _, file := range meta4.Files {
			rec := s3StemcellRec{
				Key: file.Name,
				// etag is not set
				Size: file.Size,
				URL:  file.URLs[0].URL,
			}

			if meta4.Published != nil {
				rec.LastModified = meta4.Published.UTC().Format(time.RFC3339)
			}

			for _, hash := range file.Hashes {
				if hash.Type == "sha-1" {
					rec.SHA1 = hash.Hash
				}
				if hash.Type == "sha-256" {
					rec.SHA256 = hash.Hash
				}
			}

			recs = append(recs, rec)
		}
	}

	return recs, nil
}

func (r s3StemcellRec) Validate() error {
	if len(r.Key) == 0 {
		return bosherr.Error("Expected stemcell rec to have non-empty 'Key'")
	}
	if len(r.SHA1) == 0 && len(r.ETag) == 0 {
		return bosherr.Error("Expected stemcell rec to have non-empty 'SHA1' or 'ETag'")
	}
	if r.Size == 0 {
		return bosherr.Error("Expected stemcell rec to have non-empty 'Size'")
	}
	if len(r.URL) == 0 {
		return bosherr.Error("Expected stemcell rec to have non-empty 'URL'")
	}
	return nil
}
