package jobsrepo

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	bpindex "github.com/cppforlife/bosh-provisioner/index"
	bpreljob "github.com/cppforlife/bosh-provisioner/release/job"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

type CJRepository struct {
	index  bpindex.Index
	logger boshlog.Logger
}

type relVerRecKey struct {
	Source     string
	VersionRaw string
}

func NewConcreteJobsRepository(
	index bpindex.Index,
	logger boshlog.Logger,
) CJRepository {
	return CJRepository{
		index:  index,
		logger: logger,
	}
}

func (r CJRepository) FindAll(relVerRec bhrelsrepo.ReleaseVersionRec) ([]bpreljob.Job, bool, error) {
	var relJobs []bpreljob.Job

	key := relVerRecKey{Source: relVerRec.Source, VersionRaw: relVerRec.VersionRaw}

	err := r.index.Find(key, &relJobs)
	if err != nil {
		if err == bpindex.ErrNotFound {
			return relJobs, false, nil
		}

		return relJobs, false, bosherr.WrapError(err, "Finding release jobs")
	}

	return relJobs, true, nil
}

func (r CJRepository) SaveAll(relVerRec bhrelsrepo.ReleaseVersionRec, relJobs []bpreljob.Job) error {
	key := relVerRecKey{Source: relVerRec.Source, VersionRaw: relVerRec.VersionRaw}

	err := r.index.Save(key, relJobs)
	if err != nil {
		return bosherr.WrapError(err, "Saving release jobs")
	}

	return nil
}
