package releasesrepo

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	bprel "github.com/bosh-dep-forks/bosh-provisioner/release"
	bhrelver "github.com/bosh-io/web/release/relver"
)

type CRVRepository struct {
	relVerFactory bhrelver.Factory
	logger        boshlog.Logger
}

func NewConcreteReleaseVersionsRepository(
	relVerFactory bhrelver.Factory,
	logger boshlog.Logger,
) CRVRepository {
	return CRVRepository{
		relVerFactory: relVerFactory,
		logger:        logger,
	}
}

func (r CRVRepository) Find(relVerRec ReleaseVersionRec) (bprel.Release, error) {
	var rel bprel.Release

	relVer, err := r.relVerFactory.Find(relVerRec.Source, relVerRec.VersionRaw)
	if err != nil {
		return rel, err
	}

	contents, err := relVer.Read("release.v1.yml")
	if err != nil {
		return rel, err
	}

	err = json.Unmarshal(contents, &rel)
	if err != nil {
		return rel, bosherr.WrapError(err, "Unmarshaling release")
	}

	return rel, nil
}
