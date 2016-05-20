package fetcher

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	bhtarball "github.com/cppforlife/bosh-hub/release/fetcher/tarball"
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

type TarballRelease struct {
	manifestPath string

	builder   bhtarball.Builder
	extractor bhtarball.Extractor
	uploader  bhtarball.Uploader

	releasesRepo bhrelsrepo.ReleasesRepository

	logTag string
	logger boshlog.Logger
}

func NewTarballRelease(
	manifestPath string,
	builder bhtarball.Builder,
	extractor bhtarball.Extractor,
	uploader bhtarball.Uploader,
	releasesRepo bhrelsrepo.ReleasesRepository,
	logger boshlog.Logger,
) TarballRelease {
	return TarballRelease{
		manifestPath: manifestPath,

		builder:   builder,
		extractor: extractor,
		uploader:  uploader,

		releasesRepo: releasesRepo,

		logTag: "TarballRelease",
		logger: logger,
	}
}

func (tr TarballRelease) Import(url string) error {
	tr.logger.Debug(tr.logTag, "Importing tarball release from '%s'", url)

	tgzPath, err := tr.builder.Build(tr.manifestPath)
	if err != nil {
		return bosherr.WrapError(err, "Building release tarball from manifest '%s'", tr.manifestPath)
	}

	defer tr.builder.CleanUp(tgzPath)

	relVerRec, err := tr.extractor.Extract(url, tgzPath)
	if err != nil {
		return bosherr.WrapError(err, "Importing release version from tgz '%s'", tgzPath)
	}

	err = tr.uploader.Upload(relVerRec, tgzPath)
	if err != nil {
		return bosherr.WrapError(err, "Uploading release verison '%v'", relVerRec)
	}

	// Save release version only after everything else was successfully imported
	err = tr.releasesRepo.Add(relVerRec)
	if err != nil {
		return bosherr.WrapError(err, "Saving release version into release versions repository")
	}

	return nil
}
