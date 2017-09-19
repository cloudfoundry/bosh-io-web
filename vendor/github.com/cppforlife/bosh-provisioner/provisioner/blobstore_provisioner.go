package provisioner

import (
	"os"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type BlobstoreProvisioner struct {
	fs              boshsys.FileSystem
	blobstoreConfig BlobstoreConfig
	logger          boshlog.Logger
}

func NewBlobstoreProvisioner(
	fs boshsys.FileSystem,
	blobstoreConfig BlobstoreConfig,
	logger boshlog.Logger,
) BlobstoreProvisioner {
	return BlobstoreProvisioner{
		fs:              fs,
		blobstoreConfig: blobstoreConfig,
		logger:          logger,
	}
}

func (p BlobstoreProvisioner) Provision() error {
	blobstorePath := p.blobstoreConfig.LocalPath()
	if blobstorePath != "" {
		err := p.fs.MkdirAll(blobstorePath, os.ModeDir)
		if err != nil {
			return err
		}
	}

	return nil
}
