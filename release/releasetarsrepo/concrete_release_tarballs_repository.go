package releasetarsrepo

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	bpindex "github.com/cppforlife/bosh-provisioner/index"

	bhs3 "github.com/cppforlife/bosh-hub/s3"
)

type CRTRepository struct {
	index      bpindex.Index
	urlFactory bhs3.URLFactory
	logger     boshlog.Logger
}

type releaseVersionRecKey struct {
	Source     string
	VersionRaw string
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

func (r CRTRepository) Find(source, version string) (ReleaseTarballRec, error) {
	var relTarRec ReleaseTarballRec

	key := releaseVersionRecKey{Source: source, VersionRaw: version}

	err := r.index.Find(key, &relTarRec)
	if err != nil {
		return relTarRec, bosherr.WrapError(err, "Finding release tarball")
	}

	relTarRec.urlFactory = r.urlFactory
	relTarRec.source = source
	relTarRec.versionRaw = version

	return relTarRec, nil
}

func (r CRTRepository) Save(source, version string, relTarRec ReleaseTarballRec) error {
	key := releaseVersionRecKey{Source: source, VersionRaw: version}

	err := r.index.Save(key, relTarRec)
	if err != nil {
		return bosherr.WrapError(err, "Saving release tarball")
	}

	return nil
}
