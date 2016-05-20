package repo

import (
	"sort"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	bpindex "github.com/cppforlife/bosh-provisioner/index"
	semiver "github.com/cppforlife/go-semi-semantic/version"

	bhs3 "github.com/cppforlife/bosh-hub/s3"
)

var (
	// s3BinaryKeyLatest is the single key used to keep S3 bucket results
	s3BinaryKeyLatest = s3BinaryKey{"latest"}
)

type S3Repository struct {
	index  bpindex.Index
	logger boshlog.Logger
}

type s3BinaryKey struct {
	Key string
}

type s3BinaryRec struct {
	Key  string
	ETag string

	Size         uint64
	LastModified string

	URL string
}

type VersionSorting []semiver.Version

type BinaryVersionSorting []BinaryGroup

func NewS3Repository(
	index bpindex.Index,
	logger boshlog.Logger,
) S3Repository {
	return S3Repository{
		index:  index,
		logger: logger,
	}
}

func (r S3Repository) FindLatest() ([]BinaryGroup, error) {
	binaries, err := r.findAll()
	if err != nil {
		return []BinaryGroup{}, err
	}

	versions := []semiver.Version{}

	groups := map[string][]Binary{}

	for _, b := range binaries {
		v := b.Version()

		if _, ok := groups[v.String()]; !ok {
			versions = append(versions, v)
		}

		groups[v.String()] = append(groups[v.String()], b)
	}

	sort.Sort(sort.Reverse(VersionSorting(versions)))

	latestBinGroups := []BinaryGroup{}

	for i, v := range versions {
		if i > 0 {
			break
		}

		binGroup := BinaryGroup{
			Version:  v,
			Binaries: groups[v.String()],
		}

		latestBinGroups = append(latestBinGroups, binGroup)
	}

	return latestBinGroups, nil
}

func (r S3Repository) findAll() ([]Binary, error) {
	var binaries []Binary

	var s3BinaryRecs []s3BinaryRec

	err := r.index.Find(s3BinaryKeyLatest, &s3BinaryRecs)
	if err != nil {
		if err == bpindex.ErrNotFound {
			return binaries, nil
		}

		return binaries, bosherr.WrapError(err, "Finding all S3 binary file recs")
	}

	for _, rec := range s3BinaryRecs {
		binary := NewS3Binary(
			rec.Key,
			rec.ETag,
			rec.Size,
			rec.LastModified,
			rec.URL,
		)
		if binary == nil {
			continue
		}

		binaries = append(binaries, binary)
	}

	return binaries, nil
}

func (r S3Repository) SaveAll(s3Files []bhs3.File) error {
	var s3BinaryRecs []s3BinaryRec

	for _, s3File := range s3Files {
		url, err := s3File.URL()
		if err != nil {
			return bosherr.WrapError(err, "Generating S3 binary file URL")
		}

		rec := s3BinaryRec{
			Key:  s3File.Key(),
			ETag: s3File.ETag(),

			Size:         s3File.Size(),
			LastModified: s3File.LastModified(),

			URL: url,
		}

		s3BinaryRecs = append(s3BinaryRecs, rec)
	}

	err := r.index.Save(s3BinaryKeyLatest, s3BinaryRecs)
	if err != nil {
		return bosherr.WrapError(err, "Saving all S3 binary file recs")
	}

	return nil
}

func (s VersionSorting) Len() int           { return len(s) }
func (s VersionSorting) Less(i, j int) bool { return s[i].IsLt(s[j]) }
func (s VersionSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s BinaryVersionSorting) Len() int           { return len(s) }
func (s BinaryVersionSorting) Less(i, j int) bool { return s[i].Version.IsLt(s[j].Version) }
func (s BinaryVersionSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
