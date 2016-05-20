package importerrsrepo

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	bpindex "github.com/cppforlife/bosh-provisioner/index"

	bhimpsrepo "github.com/cppforlife/bosh-hub/release/importsrepo"
)

type CIERepository struct {
	index  bpindex.Index
	logTag string
	logger boshlog.Logger
}

type importErrRecKey struct {
	RelSource string
	Version   string
}

func NewConcreteImportErrsRepository(index bpindex.Index, logger boshlog.Logger) CIERepository {
	return CIERepository{
		index:  index,
		logTag: "ConcreteImportErrsRepository",
		logger: logger,
	}
}

func (r CIERepository) ListAll() ([]ImportErrRec, error) {
	var importErrRecs []ImportErrRec

	err := r.index.List(&importErrRecs)
	if err != nil {
		return nil, bosherr.WrapError(err, "Finding import errs")
	}

	return importErrRecs, nil
}

func (r CIERepository) Add(importErrRec ImportErrRec) error {
	err := importErrRec.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating import err")
	}

	key := importErrRecKey{
		RelSource: importErrRec.ImportRec.RelSource,
		Version:   importErrRec.ImportRec.Version,
	}

	// Since key is unique it is ok to dup the item; value is overriden
	err = r.index.Save(key, importErrRec)
	if err != nil {
		return bosherr.WrapError(err, "Saving import err")
	}

	return nil
}

func (r CIERepository) Remove(importRec bhimpsrepo.ImportRec) error {
	key := importErrRecKey{
		RelSource: importRec.RelSource,
		Version:   importRec.Version,
	}

	err := r.index.Remove(key)
	if err != nil {
		return bosherr.WrapError(err, "Removing import err")
	}

	return nil
}

func (r CIERepository) Contains(importRec bhimpsrepo.ImportRec) (bool, error) {
	key := importErrRecKey{
		RelSource: importRec.RelSource,
		Version:   importRec.Version,
	}

	err := r.index.Find(key, nil)
	if err != nil {
		if err == bpindex.ErrNotFound {
			return false, nil
		}

		return false, bosherr.WrapError(err, "Finding import err")
	}

	return true, nil
}
