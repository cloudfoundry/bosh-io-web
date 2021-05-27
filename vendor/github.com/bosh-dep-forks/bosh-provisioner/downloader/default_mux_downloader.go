package downloader

import (
	boshblob "github.com/cloudfoundry/bosh-utils/blobstore"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

func NewDefaultMuxDownloader(
	fs boshsys.FileSystem,
	runner boshsys.CmdRunner,
	blobstore boshblob.DigestBlobstore,
	logger boshlog.Logger,
) MuxDownloader {
	mux := map[string]Downloader{
		"http":  NewHTTPDownloader(fs, logger),
		"https": NewHTTPDownloader(fs, logger),
		"file":  NewLocalFSDownloader(fs, logger),
	}

	if runner != nil {
		mux["git"] = NewGitDownloader(fs, runner, logger)
	}

	if blobstore != nil {
		mux["blobstore"] = NewBlobstoreDownloader(blobstore, logger)
	}

	return NewMuxDownloader(mux, logger)
}
