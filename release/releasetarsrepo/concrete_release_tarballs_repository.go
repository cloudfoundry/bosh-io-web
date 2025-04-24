package releasetarsrepo

import (
	"encoding/xml"
	"errors"
	"path/filepath"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	bhrelver "github.com/bosh-io/web/release/relver"
	bhs3 "github.com/bosh-io/web/s3"

	"github.com/dpb587/metalink"
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

	var meta4 metalink.Metalink

	err = xml.Unmarshal(contents, &meta4)
	if err != nil {
		return relTarRec, bosherr.WrapError(err, "Unmarshaling meta4")
	}

	relTarRec.BlobID = filepath.Base(meta4.Files[0].URLs[0].URL)

	for _, hash := range meta4.Files[0].Hashes {
		if hash.Type == "sha-1" {
			relTarRec.SHA1 = hash.Hash
		}
		if hash.Type == "sha-256" {
			relTarRec.SHA256 = hash.Hash
		}
	}

	if relTarRec.SHA1 == "" {
		return relTarRec, errors.New("Missing SHA1") //nolint:staticcheck
	}

	return relTarRec, nil
}
