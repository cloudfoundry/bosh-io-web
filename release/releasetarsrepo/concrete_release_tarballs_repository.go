package releasetarsrepo

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	bpindex "github.com/cppforlife/bosh-provisioner/index"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhs3 "github.com/cppforlife/bosh-hub/s3"
)

type CRTRepository struct {
	index      bpindex.Index
	urlFactory bhs3.URLFactory
	logger     boshlog.Logger
}

func NewConcreteReleaseTarballsRepository(
	index bpindex.Index,
	urlFactory bhs3.URLFactory,
	logger boshlog.Logger,
) CRTRepository {
	return CRTRepository{
		index:      index,
		urlFactory: urlFactory,
		logger:     logger,
	}
}

func (r CRTRepository) Find(relVerRec bhrelsrepo.ReleaseVersionRec) (ReleaseTarballRec, bool, error) {
	var relTarRec ReleaseTarballRec

	err := r.index.Find(relVerRec, &relTarRec)
	if err != nil {
		if err == bpindex.ErrNotFound {
			return relTarRec, false, nil
		}

		return relTarRec, false, bosherr.WrapError(err, "Finding release tarball")
	}

	relTarRec.urlFactory = r.urlFactory
	relTarRec.relVerRec = relVerRec

	return relTarRec, true, nil
}

func (r CRTRepository) Save(relVerRec bhrelsrepo.ReleaseVersionRec, relTarRec ReleaseTarballRec) error {
	err := r.index.Save(relVerRec, relTarRec)
	if err != nil {
		return bosherr.WrapError(err, "Saving release tarball")
	}

	return nil
}
