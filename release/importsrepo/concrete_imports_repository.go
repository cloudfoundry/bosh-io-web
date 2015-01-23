package importsrepo

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	bpindex "github.com/cppforlife/bosh-provisioner/index"
)

const ciRepositoryLogTag = "ConcreteImportsRepository"

type CIRepository struct {
	index  bpindex.Index
	logger boshlog.Logger
}

type relSourceRec struct {
	RelSource string
	Version   string
}

func NewConcreteImportsRepository(index bpindex.Index, logger boshlog.Logger) CIRepository {
	return CIRepository{
		index:  index,
		logger: logger,
	}
}

func (r CIRepository) ListAll() ([]ImportRec, error) {
	var importRecs []ImportRec
	var relSourceRecs []relSourceRec

	err := r.index.ListKeys(&relSourceRecs)
	if err != nil {
		return nil, bosherr.WrapError(err, "Finding release sources")
	}

	for _, relSourceRec := range relSourceRecs {
		importRec := ImportRec{
			RelSource: relSourceRec.RelSource,
			Version:   relSourceRec.Version,
		}

		importRecs = append(importRecs, importRec)
	}

	return importRecs, nil
}

func (r CIRepository) Push(relSource, version string) error {
	r.logger.Debug(ciRepositoryLogTag,
		"Pushing release source '%s' of version '%s'", relSource, version)

	// Since key is unique it is ok to dup the item; value is dummy
	err := r.index.Save(relSourceRec{relSource, version}, "dummy")
	if err != nil {
		return bosherr.WrapError(err, "Saving release source")
	}

	return nil
}

func (r CIRepository) Pull() (ImportRec, bool, error) {
	var relSourceRecs []relSourceRec

	err := r.index.ListKeys(&relSourceRecs)
	if err != nil {
		return ImportRec{}, false, bosherr.WrapError(err, "Finding release sources")
	}

	if len(relSourceRecs) == 0 {
		return ImportRec{}, false, nil
	}

	rec := relSourceRecs[0]

	r.logger.Debug(ciRepositoryLogTag, "Pulling import '%v'", rec)

	importRec := ImportRec{
		RelSource: rec.RelSource,
		Version:   rec.Version,
	}

	err = r.index.Remove(relSourceRec{rec.RelSource, rec.Version})
	if err != nil {
		return ImportRec{}, false, bosherr.WrapError(err, "Removing import")
	}

	return importRec, true, nil
}

func (r CIRepository) Remove(importRec ImportRec) error {
	err := r.index.Remove(relSourceRec{importRec.RelSource, importRec.Version})
	if err != nil {
		return bosherr.WrapError(err, "Removing import")
	}

	return nil
}
