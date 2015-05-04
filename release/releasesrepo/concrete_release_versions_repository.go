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

type relVerRecKey struct {
	Source     string
	VersionRaw string
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

	key := relVerRecKey{Source: relVerRec.Source, VersionRaw: relVerRec.VersionRaw}

	err := r.index.Find(key, &rel)
	if err != nil {
		if err == bpindex.ErrNotFound {
			return rel, false, nil
		}

		return rel, false, bosherr.WrapError(err, "Finding release")
	}

	return rel, true, nil
}

func (r CRVRepository) Save(relVerRec ReleaseVersionRec, rel bprel.Release) error {
	key := relVerRecKey{Source: relVerRec.Source, VersionRaw: relVerRec.VersionRaw}

	err := r.index.Save(key, rel)
	if err != nil {
		return bosherr.WrapError(err, "Saving release")
	}

	return nil
}
