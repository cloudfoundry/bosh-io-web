package fetcher

import (
	"path/filepath"

	boshblob "github.com/cloudfoundry/bosh-utils/blobstore"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"
)

type BlobstoreProvider struct {
	fs      boshsys.FileSystem
	runner  boshsys.CmdRunner
	uuidGen boshuuid.Generator
	logger  boshlog.Logger
}

func NewBlobstoreProvider(
	fs boshsys.FileSystem,
	runner boshsys.CmdRunner,
	uuidGen boshuuid.Generator,
	logger boshlog.Logger,
) BlobstoreProvider {
	return BlobstoreProvider{
		fs:      fs,
		runner:  runner,
		uuidGen: uuidGen,
		logger:  logger,
	}
}

func (p BlobstoreProvider) Get(provider string, options map[string]interface{}) (boshblob.Blobstore, error) {
	configDir, err := p.fs.TempDir("blobstore-s3-config")
	if err != nil {
		return nil, bosherr.WrapError(err, "Cerating tmp dir for blobstore config")
	}

	configPath := filepath.Join(configDir, "config.json")

	blobstore := boshblob.NewExternalBlobstore(
		provider,
		options,
		p.fs,
		p.runner,
		p.uuidGen,
		configPath,
	)

	blobstore = boshblob.NewSHA1VerifiableBlobstore(blobstore)

	blobstore = boshblob.NewRetryableBlobstore(blobstore, 3, p.logger)

	err = blobstore.Validate()
	if err != nil {
		return nil, bosherr.WrapError(err, "Validating blobstore")
	}

	return blobstore, nil
}
