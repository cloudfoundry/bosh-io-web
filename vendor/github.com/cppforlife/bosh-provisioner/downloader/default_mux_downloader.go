package downloader

import (
	boshblob "github.com/cloudfoundry/bosh-utils/blobstore"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

func NewDefaultMuxDownloader(
	fs boshsys.FileSystem,
	runner boshsys.CmdRunner,
	blobstore boshblob.Blobstore,
	logger boshlog.Logger,
) MuxDownloader {
	mux := map[string]Downloader{
		"file":  NewLocalFSDownloader(fs, logger),
	}
	return NewMuxDownloader(mux, logger)
}
