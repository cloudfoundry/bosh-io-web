package releasetarsrepo

import (
	"encoding/xml"
	"path/filepath"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	bhrelver "github.com/cppforlife/bosh-hub/release/relver"
	bhs3 "github.com/cppforlife/bosh-hub/s3"
)

type ReleaseTarballsRepository interface {
	Find(source, version string) (ReleaseTarballRec, error)
}

type CRTRepository struct {
	relVerFactory bhrelver.Factory
	urlFactory    bhs3.URLFactory
	logger        boshlog.Logger
}

func NewConcreteReleaseTarballsRepository(
	relVerFactory bhrelver.Factory,
	urlFactory bhs3.URLFactory,
	logger boshlog.Logger,
) CRTRepository {
	return CRTRepository{
		relVerFactory: relVerFactory,
		urlFactory:    urlFactory,
		logger:        logger,
	}
}

func (r CRTRepository) Find(source, version string) (ReleaseTarballRec, error) {
	var relTarRec ReleaseTarballRec

	relVer, err := r.relVerFactory.Find(source, version)
	if err != nil {
		return relTarRec, err
	}

	contents, err := relVer.Read("source.meta4")
	if err != nil {
		return relTarRec, err
	}

	relTarRec.urlFactory = r.urlFactory
	relTarRec.source = source
	relTarRec.versionRaw = version

	var meta4 Metalink

	err = xml.Unmarshal(contents, &meta4)
	if err != nil {
		return relTarRec, bosherr.WrapError(err, "Unmarshaling meta4")
	}

	relTarRec.BlobID = filepath.Base(meta4.Files[0].URLs[0].URL)
	relTarRec.SHA1 = meta4.Files[0].Hashes[0].Hash

	return relTarRec, nil
}
