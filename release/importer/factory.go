package importer

import (
	"time"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	boshuuid "github.com/cloudfoundry/bosh-agent/uuid"
	bpdload "github.com/cppforlife/bosh-provisioner/downloader"
	bprel "github.com/cppforlife/bosh-provisioner/release"
	bpreljob "github.com/cppforlife/bosh-provisioner/release/job"
	bptar "github.com/cppforlife/bosh-provisioner/tar"

	bhfetcher "github.com/cppforlife/bosh-hub/release/fetcher"
	bhimperrsrepo "github.com/cppforlife/bosh-hub/release/importerrsrepo"
	bhimpsrepo "github.com/cppforlife/bosh-hub/release/importsrepo"
	bhjobsrepo "github.com/cppforlife/bosh-hub/release/jobsrepo"
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

type FactoryOptions struct {
	Enabled bool
	Period  time.Duration

	BlobstoreType    string
	BlobstoreOptions map[string]interface{}
}

type FactoryRepos interface {
	ReleasesRepo() bhrelsrepo.ReleasesRepository
	ReleaseVersionsRepo() bhrelsrepo.ReleaseVersionsRepository
	JobsRepo() bhjobsrepo.JobsRepository
	ImportsRepo() bhimpsrepo.ImportsRepository
	ImportErrsRepo() bhimperrsrepo.ImportErrsRepository
}

type Factory struct {
	Importer Importer
}

func NewFactory(
	options FactoryOptions,
	repos FactoryRepos,
	fetcher bhfetcher.Fetcher,
	fs boshsys.FileSystem,
	runner boshsys.CmdRunner,
	downloader bpdload.Downloader,
	uuidGen boshuuid.Generator,
	logger boshlog.Logger,
) (Factory, error) {
	if !options.Enabled {
		return Factory{Importer: NewNoopImporter(logger)}, nil
	}

	var tarballReleaseFactory bhfetcher.TarballReleaseFactory

	{
		blobstoreProvider := bhfetcher.NewBlobstoreProvider(fs, runner, uuidGen, logger)

		blobstore, err := blobstoreProvider.Get(options.BlobstoreType, options.BlobstoreOptions)
		if err != nil {
			return Factory{}, bosherr.WrapError(err, "Building blobstore")
		}

		extractor := bptar.NewCmdExtractor(runner, fs, logger)

		releaseReaderFactory := bprel.NewReaderFactory(downloader, extractor, fs, logger)

		jobReaderFactory := bpreljob.NewReaderFactory(downloader, extractor, fs, logger)

		tarballReleaseFactory = bhfetcher.NewTarballReleaseFactory(
			fs,
			runner,
			blobstore,
			releaseReaderFactory,
			jobReaderFactory,
			repos.ReleasesRepo(),
			repos.ReleaseVersionsRepo(),
			repos.JobsRepo(),
			logger,
		)
	}

	periodicImporter := NewQueueImporter(
		options.Period,
		make(chan struct{}),
		tarballReleaseFactory,
		repos.ReleasesRepo(),
		repos.ImportsRepo(),
		repos.ImportErrsRepo(),
		fetcher,
		logger,
	)

	return Factory{Importer: periodicImporter}, nil
}
