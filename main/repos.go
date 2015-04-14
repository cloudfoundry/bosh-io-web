package main

import (
	"path/filepath"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	bpindex "github.com/cppforlife/bosh-provisioner/index"

	bhindex "github.com/cppforlife/bosh-hub/index"
	bhimperrsrepo "github.com/cppforlife/bosh-hub/release/importerrsrepo"
	bhimpsrepo "github.com/cppforlife/bosh-hub/release/importsrepo"
	bhjobsrepo "github.com/cppforlife/bosh-hub/release/jobsrepo"
	bhnotesrepo "github.com/cppforlife/bosh-hub/release/notesrepo"
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhreltarsrepo "github.com/cppforlife/bosh-hub/release/releasetarsrepo"
	bhwatchersrepo "github.com/cppforlife/bosh-hub/release/watchersrepo"
	bhs3 "github.com/cppforlife/bosh-hub/s3"
	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
)

type ReposOptions struct {
	Type string // e.g. file, db

	Dir     string
	ConnURL string

	PredefinedReleaseSources []string
	PredefinedAvatars        map[string]string

	ReleaseTarballLinker ReleaseTarballLinkerOptions
}

type ReleaseTarballLinkerOptions struct {
	Type string // e.g. CloudFront, S3

	BaseURL string

	KeyPairID  string
	PrivateKey string
}

type Repos struct {
	releasesRepo        bhrelsrepo.ReleasesRepository
	releaseTarsRepo     bhreltarsrepo.ReleaseTarballsRepository
	releaseVersionsRepo bhrelsrepo.ReleaseVersionsRepository
	jobsRepo            bhjobsrepo.JobsRepository

	s3StemcellsRepo bhstemsrepo.S3StemcellsRepository

	importsRepo    bhimpsrepo.ImportsRepository
	importErrsRepo bhimperrsrepo.ImportErrsRepository
	watchersRepo   bhwatchersrepo.WatchersRepository
}

func NewRepos(options ReposOptions, fs boshsys.FileSystem, logger boshlog.Logger) (Repos, error) {
	var i repoIndicies

	var err error

	switch {
	case options.Type == "file":
		i = newFileRepoIndicies(options.Dir, fs)
	case options.Type == "db":
		i, err = newDBRepoIndicies(options.ConnURL, logger)
	default:
		err = bosherr.New("Expected repos type '%s'", options.Type)
	}
	if err != nil {
		return Repos{}, err
	}

	predefinedAvatars := options.PredefinedAvatars
	if predefinedAvatars == nil {
		predefinedAvatars = map[string]string{}
	}

	linkerOpts := options.ReleaseTarballLinker

	var linkerFactory bhs3.URLFactory

	switch {
	case linkerOpts.Type == "CloudFront":
		linkerFactory, err = bhs3.NewCDNURLFactory(linkerOpts.BaseURL, linkerOpts.KeyPairID, linkerOpts.PrivateKey)
	case linkerOpts.Type == "S3":
		linkerFactory = bhs3.NewDirectURLFactory(linkerOpts.BaseURL)
	default:
		err = bosherr.New("Expected linker type '%s'", linkerOpts.Type)
	}
	if err != nil {
		return Repos{}, err
	}

	releaseNotesRepo := bhnotesrepo.NewConcreteNotesRepository(i.releaseNotesIndex, logger)

	releasesRepo := bhrelsrepo.NewConcreteReleasesRepository(
		options.PredefinedReleaseSources,
		predefinedAvatars,
		i.releasesIndex,
		releaseNotesRepo,
		logger,
	)

	repos := Repos{
		releasesRepo: releasesRepo,

		releaseTarsRepo:     bhreltarsrepo.NewConcreteReleaseTarballsRepository(i.releaseTarsIndex, linkerFactory, logger),
		releaseVersionsRepo: bhrelsrepo.NewConcreteReleaseVersionsRepository(i.releaseVersionsIndex, logger),
		jobsRepo:            bhjobsrepo.NewConcreteJobsRepository(i.jobsIndex, logger),

		s3StemcellsRepo: bhstemsrepo.NewS3StemcellsRepository(i.s3StemcellsIndex, logger),

		importsRepo:    bhimpsrepo.NewConcreteImportsRepository(i.importsIndex, logger),
		importErrsRepo: bhimperrsrepo.NewConcreteImportErrsRepository(i.importErrsIndex, logger),
		watchersRepo:   bhwatchersrepo.NewConcreteWatchersRepository(i.watchersIndex, logger),
	}

	return repos, nil
}

func (r Repos) ReleasesRepo() bhrelsrepo.ReleasesRepository { return r.releasesRepo }

func (r Repos) ReleaseTarsRepo() bhreltarsrepo.ReleaseTarballsRepository { return r.releaseTarsRepo }

func (r Repos) ReleaseVersionsRepo() bhrelsrepo.ReleaseVersionsRepository {
	return r.releaseVersionsRepo
}

func (r Repos) JobsRepo() bhjobsrepo.JobsRepository { return r.jobsRepo }

func (r Repos) S3StemcellsRepo() bhstemsrepo.S3StemcellsRepository { return r.s3StemcellsRepo }

func (r Repos) StemcellsRepo() bhstemsrepo.StemcellsRepository { return r.s3StemcellsRepo }

func (r Repos) ImportsRepo() bhimpsrepo.ImportsRepository { return r.importsRepo }

func (r Repos) ImportErrsRepo() bhimperrsrepo.ImportErrsRepository { return r.importErrsRepo }

func (r Repos) WatchersRepo() bhwatchersrepo.WatchersRepository { return r.watchersRepo }

type repoIndicies struct {
	releasesIndex        bpindex.Index
	releaseNotesIndex    bpindex.Index
	releaseTarsIndex     bpindex.Index
	releaseVersionsIndex bpindex.Index
	jobsIndex            bpindex.Index

	s3StemcellsIndex bpindex.Index

	importsIndex    bpindex.Index
	importErrsIndex bpindex.Index
	watchersIndex   bpindex.Index
}

func newFileRepoIndicies(dir string, fs boshsys.FileSystem) repoIndicies {
	return repoIndicies{
		releasesIndex:        bpindex.NewFileIndex(filepath.Join(dir, "releases.json"), fs),
		releaseNotesIndex:    bpindex.NewFileIndex(filepath.Join(dir, "release_notes.json"), fs),
		releaseTarsIndex:     bpindex.NewFileIndex(filepath.Join(dir, "release_tarballs.json"), fs),
		releaseVersionsIndex: bpindex.NewFileIndex(filepath.Join(dir, "release_versions.json"), fs),
		jobsIndex:            bpindex.NewFileIndex(filepath.Join(dir, "jobs.json"), fs),

		s3StemcellsIndex: bpindex.NewFileIndex(filepath.Join(dir, "s3_stemcells.json"), fs),

		importsIndex:    bpindex.NewFileIndex(filepath.Join(dir, "imports.json"), fs),
		importErrsIndex: bpindex.NewFileIndex(filepath.Join(dir, "import_errs.json"), fs),
		watchersIndex:   bpindex.NewFileIndex(filepath.Join(dir, "watchers.json"), fs),
	}
}

func newDBRepoIndicies(url string, logger boshlog.Logger) (repoIndicies, error) {
	adapterPool, err := bhindex.NewPostgresAdapterPool(url, logger)
	if err != nil {
		return repoIndicies{}, err
	}

	releasesAdapter, err := adapterPool.NewAdapter("releases")
	if err != nil {
		return repoIndicies{}, err
	}

	releaseNotesAdapter, err := adapterPool.NewAdapter("release_notes")
	if err != nil {
		return repoIndicies{}, err
	}

	releasesTarballsAdapter, err := adapterPool.NewAdapter("release_tarballs")
	if err != nil {
		return repoIndicies{}, err
	}

	releaseVersionsAdapter, err := adapterPool.NewAdapter("release_versions")
	if err != nil {
		return repoIndicies{}, err
	}

	jobsAdapter, err := adapterPool.NewAdapter("jobs")
	if err != nil {
		return repoIndicies{}, err
	}

	s3StemcellsAdapter, err := adapterPool.NewAdapter("s3_stemcells")
	if err != nil {
		return repoIndicies{}, err
	}

	importsAdapter, err := adapterPool.NewAdapter("imports")
	if err != nil {
		return repoIndicies{}, err
	}

	importErrsAdapter, err := adapterPool.NewAdapter("import_errs")
	if err != nil {
		return repoIndicies{}, err
	}

	watchersAdapter, err := adapterPool.NewAdapter("watchers")
	if err != nil {
		return repoIndicies{}, err
	}

	indicies := repoIndicies{
		releasesIndex:        bhindex.NewDBIndex(releasesAdapter, logger),
		releaseNotesIndex:    bhindex.NewDBIndex(releaseNotesAdapter, logger),
		releaseTarsIndex:     bhindex.NewDBIndex(releasesTarballsAdapter, logger),
		releaseVersionsIndex: bhindex.NewDBIndex(releaseVersionsAdapter, logger),
		jobsIndex:            bhindex.NewDBIndex(jobsAdapter, logger),

		s3StemcellsIndex: bhindex.NewDBIndex(s3StemcellsAdapter, logger),

		importsIndex:    bhindex.NewDBIndex(importsAdapter, logger),
		importErrsIndex: bhindex.NewDBIndex(importErrsAdapter, logger),
		watchersIndex:   bhindex.NewDBIndex(watchersAdapter, logger),
	}

	return indicies, nil
}
