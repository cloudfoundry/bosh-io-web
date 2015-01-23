package releasesrepo

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	bpindex "github.com/cppforlife/bosh-provisioner/index"
	bprel "github.com/cppforlife/bosh-provisioner/release"
)

type CRVRepository struct {
	index  bpindex.Index
	logger boshlog.Logger
}

func NewConcreteReleaseVersionsRepository(
	index bpindex.Index,
	logger boshlog.Logger,
) CRVRepository {
	return CRVRepository{
		index:  index,
		logger: logger,
	}
}

func (r CRVRepository) Find(relVerRec ReleaseVersionRec) (bprel.Release, bool, error) {
	var rel bprel.Release

	err := r.index.Find(relVerRec, &rel)
	if err != nil {
		if err == bpindex.ErrNotFound {
			return rel, false, nil
		}

		return rel, false, bosherr.WrapError(err, "Finding release")
	}

	return rel, true, nil
}

func (r CRVRepository) Save(relVerRec ReleaseVersionRec, rel bprel.Release) error {
	err := r.index.Save(relVerRec, rel)
	if err != nil {
		return bosherr.WrapError(err, "Saving release")
	}

	return nil
}
