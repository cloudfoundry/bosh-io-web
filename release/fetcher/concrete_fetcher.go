package fetcher

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	bpdload "github.com/cppforlife/bosh-provisioner/downloader"
)

const concreteFetcherLogTag = "ConcreteFetcher"

type ConcreteFetcher struct {
	fs         boshsys.FileSystem
	downloader bpdload.Downloader
	logger     boshlog.Logger
}

func NewConcreteFetcher(
	fs boshsys.FileSystem,
	downloader bpdload.Downloader,
	logger boshlog.Logger,
) ConcreteFetcher {
	return ConcreteFetcher{
		fs:         fs,
		downloader: downloader,
		logger:     logger,
	}
}

func (f ConcreteFetcher) Fetch(relSource string) (ReleaseDir, error) {
	var releaseDir ReleaseDir

	f.logger.Debug(concreteFetcherLogTag, "Starting fetching release '%s'", relSource)

	// todo take identifier
	downloadPath, err := f.downloader.Download("git://" + relSource)
	if err != nil {
		f.logger.Error(concreteFetcherLogTag,
			"Failed to download release '%s': %s", relSource, err)
		return releaseDir, bosherr.WrapError(err, "Downloading release")
	}

	cleanUp := func() error {
		err := f.downloader.CleanUp(downloadPath)
		if err != nil {
			f.logger.Error(concreteFetcherLogTag,
				"Failed to clean up downloaded release '%s': %s", relSource, err)

			return bosherr.WrapError(err, "Cleaning up downloaded release")
		}

		return nil
	}

	releaseDir = NewReleaseDir(downloadPath, cleanUp, f.fs, f.logger)

	return releaseDir, nil
}
