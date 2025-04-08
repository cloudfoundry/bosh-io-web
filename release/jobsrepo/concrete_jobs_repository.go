package jobsrepo

import (
	"encoding/json"

	bpreljob "github.com/bosh-dep-forks/bosh-provisioner/release/job"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	bhrelsrepo "github.com/bosh-io/web/release/releasesrepo"
	bhrelver "github.com/bosh-io/web/release/relver"
)

type JobsRepository interface {
	FindAll(bhrelsrepo.ReleaseVersionRec) ([]bpreljob.Job, error)
}

type CJRepository struct {
	relVerFactory bhrelver.Factory
	logger        boshlog.Logger
}

type relVerRecKey struct { //nolint:unused
	Source     string
	VersionRaw string
}

func NewConcreteJobsRepository(
	relVerFactory bhrelver.Factory,
	logger boshlog.Logger,
) CJRepository {
	return CJRepository{
		relVerFactory: relVerFactory,
		logger:        logger,
	}
}

func (r CJRepository) FindAll(relVerRec bhrelsrepo.ReleaseVersionRec) ([]bpreljob.Job, error) {
	var relJobs []bpreljob.Job

	relVer, err := r.relVerFactory.Find(relVerRec.Source, relVerRec.VersionRaw)
	if err != nil {
		return nil, err
	}

	contents, err := relVer.Read("jobs.v1.yml")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, &relJobs)
	if err != nil {
		return nil, bosherr.WrapError(err, "Unmarshaling release jobs")
	}

	return relJobs, nil
}
