package fetcher

import (
	boshblob "github.com/cloudfoundry/bosh-agent/blobstore"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	bprel "github.com/cppforlife/bosh-provisioner/release"
	bpreljob "github.com/cppforlife/bosh-provisioner/release/job"

	bhtarball "github.com/cppforlife/bosh-hub/release/fetcher/tarball"
	bhjobsrepo "github.com/cppforlife/bosh-hub/release/jobsrepo"
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

type TarballReleaseFactory struct {
	fs        boshsys.FileSystem
	runner    boshsys.CmdRunner
	blobstore boshblob.Blobstore

	releaseReaderFactory bprel.ReaderFactory
	jobReaderFactory     bpreljob.ReaderFactory

	releasesRepo        bhrelsrepo.ReleasesRepository
	releaseVersionsRepo bhrelsrepo.ReleaseVersionsRepository
	jobsRepo            bhjobsrepo.JobsRepository

	logger boshlog.Logger
}

func NewTarballReleaseFactory(
	fs boshsys.FileSystem,
	runner boshsys.CmdRunner,
	blobstore boshblob.Blobstore,
	releaseReaderFactory bprel.ReaderFactory,
	jobReaderFactory bpreljob.ReaderFactory,
	releasesRepo bhrelsrepo.ReleasesRepository,
	releaseVersionsRepo bhrelsrepo.ReleaseVersionsRepository,
	jobsRepo bhjobsrepo.JobsRepository,
	logger boshlog.Logger,
) TarballReleaseFactory {
	return TarballReleaseFactory{
		fs:        fs,
		runner:    runner,
		blobstore: blobstore,

		releaseReaderFactory: releaseReaderFactory,
		jobReaderFactory:     jobReaderFactory,

		releasesRepo:        releasesRepo,
		releaseVersionsRepo: releaseVersionsRepo,
		jobsRepo:            jobsRepo,

		logger: logger,
	}
}

func (f TarballReleaseFactory) NewTarballRelease(manifestPath string) TarballRelease {
	builder := bhtarball.NewBuilder(f.fs, f.runner, f.logger)

	extractor := bhtarball.NewExtractor(
		f.releaseReaderFactory,
		f.jobReaderFactory,
		f.releasesRepo,
		f.releaseVersionsRepo,
		f.jobsRepo,
		f.logger,
	)

	uploader := bhtarball.NewUploader(
		f.blobstore,
		f.logger,
	)

	tarballRelease := NewTarballRelease(
		manifestPath,
		builder,
		extractor,
		uploader,
		f.releasesRepo,
		f.logger,
	)

	return tarballRelease
}
