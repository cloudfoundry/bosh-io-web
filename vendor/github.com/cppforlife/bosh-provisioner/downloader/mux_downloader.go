package downloader

import (
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

const muxDownloaderLogTag = "Downloader"

type MuxDownloader struct {
	// e.g. {"http": NewHTTPDownloader()}
	mux map[string]Downloader

	logger boshlog.Logger

	// Track which Downloader should be used to clean up
	downloadedPaths map[string]Downloader
}

func NewMuxDownloader(
	mux map[string]Downloader,
	logger boshlog.Logger,
) MuxDownloader {
	return MuxDownloader{
		mux:    mux,
		logger: logger,

		downloadedPaths: map[string]Downloader{},
	}
}

func (d MuxDownloader) Download(url string) (string, error) {
	for prefix, downloader := range d.mux {
		if strings.HasPrefix(url, prefix+"://") {
			path, err := downloader.Download(url)
			if err != nil {
				return path, err
			}

			// Only remember path if Downloader succeeded
			d.downloadedPaths[path] = downloader

			return path, err
		}
	}

	return "", bosherr.Errorf("URL %s without matching downloader", url)
}

func (d MuxDownloader) CleanUp(path string) error {
	downloader, ok := d.downloadedPaths[path]
	if !ok {
		// programmer error
		return bosherr.Errorf("Unknown path %s requested to be cleaned up", path)
	}

	err := downloader.CleanUp(path)
	if err != nil {
		return err
	}

	// Forget path only if associated Downloader succeeded cleaning up
	// so that CleanUp could be called multiple times
	delete(d.downloadedPaths, path)

	return nil
}
