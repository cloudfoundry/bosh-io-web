package checksumsrepo

import (
	"reflect"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	bpindex "github.com/cppforlife/bosh-provisioner/index"
)

type CCRepository struct {
	index  bpindex.Index
	logger boshlog.Logger
}

type checksumRecKey struct {
	Key string
}

func (r checksumRecKey) Validate() error {
	if len(r.Key) == 0 {
		return bosherr.New("Expected checksum key to be non-empty")
	}

	return nil
}

func NewConcreteChecksumsRepository(
	index bpindex.Index,
	logger boshlog.Logger,
) CCRepository {
	return CCRepository{
		index:  index,
		logger: logger,
	}
}

func (r CCRepository) Find(key string) (ChecksumRec, error) {
	var rec ChecksumRec

	recKey := checksumRecKey{Key: key}

	err := recKey.Validate()
	if err != nil {
		return rec, err
	}

	err = r.index.Find(recKey, &rec)
	if err != nil {
		return rec, bosherr.WrapError(err, "Finding checksum")
	}

	return rec, nil
}

// Save uses insert to disallow checksum information to be overwritten
// (succeeds when same checksum information is provided)
func (r CCRepository) Save(key string, rec ChecksumRec) error {
	recKey := checksumRecKey{Key: key}

	err := recKey.Validate()
	if err != nil {
		return err
	}

	err = rec.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating checksum")
	}

	// Succeeds if key already holds same checksum information
	existingRec, err := r.Find(key)
	if err == nil && reflect.DeepEqual(existingRec, rec) {
		return nil
	}

	err = r.index.Insert(recKey, rec)
	if err != nil {
		return bosherr.WrapError(err, "Inserting checksum")
	}

	return nil
}
