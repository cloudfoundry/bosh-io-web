package main

import (
	"path/filepath"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	bpindex "github.com/cppforlife/bosh-provisioner/index"

	bhchecksumsrepo "github.com/cppforlife/bosh-hub/checksumsrepo"
	bhindex "github.com/cppforlife/bosh-hub/index"
	bhrelver "github.com/cppforlife/bosh-hub/release/relver"
	bhjobsrepo "github.com/cppforlife/bosh-hub/release/jobsrepo"
	bhnotesrepo "github.com/cppforlife/bosh-hub/release/notesrepo"
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhreltarsrepo "github.com/cppforlife/bosh-hub/release/releasetarsrepo"
	bhs3 "github.com/cppforlife/bosh-hub/s3"
	bhstemnotesrepo "github.com/cppforlife/bosh-hub/stemcell/notesrepo"
	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
)

type ReposOptions struct {
	Type string // e.g. file, db

	Dir     string
	ConnURL string

	ReleasesDir string
	ReleasesIndexDir string

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
	releaseVersionsRepo bhrelsrepo.ReleaseVersionsRepository
	jobsRepo            bhjobsrepo.JobsRepository

	s3StemcellsRepo    bhstemsrepo.S3StemcellsRepository
	stemcellNotesRepo  bhstemnotesrepo.NotesRepository

	checksumsRepo bhchecksumsrepo.ChecksumsRepository
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

	relVerFactory := bhrelver.NewFactory(options.ReleasesIndexDir, fs, logger)

	releaseNotesRepo := bhnotesrepo.NewConcreteNotesRepository(i.releaseNotesIndex, logger)
	releaseTarsRepo := bhreltarsrepo.NewConcreteReleaseTarballsRepository(
		options.ReleasesIndexDir, linkerFactory, fs, logger)

	releasesRepo := bhrelsrepo.NewConcreteReleasesRepository(
		options.ReleasesDir,
		options.ReleasesIndexDir,
		releaseNotesRepo,
		releaseTarsRepo,
		fs,
		logger,
	)

	checksumsRepo := bhchecksumsrepo.NewConcreteChecksumsRepository(i.checksumsIndex, logger)
	stemcellNotesRepo := bhstemnotesrepo.NewConcreteNotesRepository(i.stemcellNotesIndex, logger)

	repos := Repos{
		releasesRepo: releasesRepo,

		releaseVersionsRepo: bhrelsrepo.NewConcreteReleaseVersionsRepository(options.ReleasesIndexDir, fs, logger),
		jobsRepo:            bhjobsrepo.NewConcreteJobsRepository(relVerFactory, logger),

		s3StemcellsRepo:    bhstemsrepo.NewS3StemcellsRepository(i.s3StemcellsIndex, checksumsRepo, stemcellNotesRepo, logger),
		stemcellNotesRepo:  stemcellNotesRepo,

		checksumsRepo: checksumsRepo,
	}

	return repos, nil
}

func (r Repos) ReleasesRepo() bhrelsrepo.ReleasesRepository { return r.releasesRepo }

func (r Repos) ReleaseVersionsRepo() bhrelsrepo.ReleaseVersionsRepository {
	return r.releaseVersionsRepo
}

func (r Repos) JobsRepo() bhjobsrepo.JobsRepository { return r.jobsRepo }

func (r Repos) S3StemcellsRepo() bhstemsrepo.S3StemcellsRepository { return r.s3StemcellsRepo }
func (r Repos) StemcellNotesRepo() bhstemnotesrepo.NotesRepository { return r.stemcellNotesRepo }
func (r Repos) StemcellsRepo() bhstemsrepo.StemcellsRepository     { return r.s3StemcellsRepo }

func (r Repos) ChecksumsRepo() bhchecksumsrepo.ChecksumsRepository { return r.checksumsRepo }

type repoIndicies struct {
	releaseNotesIndex    bpindex.Index
	releaseTarsIndex     bpindex.Index

	s3StemcellsIndex    bpindex.Index
	stemcellNotesIndex  bpindex.Index
	s3BoshInitBinsIndex bpindex.Index

	checksumsIndex bpindex.Index
}

func newFileRepoIndicies(dir string, fs boshsys.FileSystem) repoIndicies {
	return repoIndicies{
		releaseNotesIndex:    bpindex.NewFileIndex(filepath.Join(dir, "release_notes.json"), fs),
		releaseTarsIndex:     bpindex.NewFileIndex(filepath.Join(dir, "release_tarballs.json"), fs),

		s3StemcellsIndex:    bpindex.NewFileIndex(filepath.Join(dir, "s3_stemcells.json"), fs),
		stemcellNotesIndex:  bpindex.NewFileIndex(filepath.Join(dir, "stemcell_notes.json"), fs),
		s3BoshInitBinsIndex: bpindex.NewFileIndex(filepath.Join(dir, "s3_bosh_init_bins.json"), fs),

		checksumsIndex: bpindex.NewFileIndex(filepath.Join(dir, "checksums.json"), fs),
	}
}

func newDBRepoIndicies(url string, logger boshlog.Logger) (repoIndicies, error) {
	adapterPool, err := bhindex.NewPostgresAdapterPool(url, logger)
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

	s3StemcellsAdapter, err := adapterPool.NewAdapter("s3_stemcells")
	if err != nil {
		return repoIndicies{}, err
	}

	stemcellNotesAdapter, err := adapterPool.NewAdapter("stemcell_notes")
	if err != nil {
		return repoIndicies{}, err
	}

	checksumsAdapter, err := adapterPool.NewAdapter("checksums")
	if err != nil {
		return repoIndicies{}, err
	}

	indicies := repoIndicies{
		releaseNotesIndex:    bhindex.NewDBIndex(releaseNotesAdapter, logger),
		releaseTarsIndex:     bhindex.NewDBIndex(releasesTarballsAdapter, logger),

		s3StemcellsIndex:    bhindex.NewDBIndex(s3StemcellsAdapter, logger),
		stemcellNotesIndex:  bhindex.NewDBIndex(stemcellNotesAdapter, logger),

		checksumsIndex: bhindex.NewDBIndex(checksumsAdapter, logger),
	}

	return indicies, nil
}
