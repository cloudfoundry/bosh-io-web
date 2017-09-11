package main

import (
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"

	bhjobsrepo "github.com/cppforlife/bosh-hub/release/jobsrepo"
	bhnotesrepo "github.com/cppforlife/bosh-hub/release/notesrepo"
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhreltarsrepo "github.com/cppforlife/bosh-hub/release/releasetarsrepo"
	bhrelver "github.com/cppforlife/bosh-hub/release/relver"
	bhs3 "github.com/cppforlife/bosh-hub/s3"
	bhstemnotesrepo "github.com/cppforlife/bosh-hub/stemcell/notesrepo"
	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
)

type ReposOptions struct {
	Type string // e.g. file, db

	Dir     string
	ConnURL string

	ReleasesDir       string
	ReleasesIndexDir  string
	StemcellsIndexDir string

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

	s3StemcellsRepo   bhstemsrepo.S3StemcellsRepository
	stemcellNotesRepo bhstemnotesrepo.NotesRepository
}

func NewRepos(options ReposOptions, fs boshsys.FileSystem, logger boshlog.Logger) (Repos, error) {
	fs = NewCachingFileSystem(fs, logger)

	var err error

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
	releaseNotesRepo := bhnotesrepo.NewConcreteNotesRepository(relVerFactory, logger)
	releaseTarsRepo := bhreltarsrepo.NewConcreteReleaseTarballsRepository(relVerFactory, linkerFactory, logger)

	releasesRepo := bhrelsrepo.NewConcreteReleasesRepository(
		options.ReleasesDir,
		options.ReleasesIndexDir,
		releaseNotesRepo,
		releaseTarsRepo,
		fs,
		logger,
	)

	stemcellNotesRepo := bhstemnotesrepo.NewConcreteNotesRepository(options.StemcellsIndexDir, fs, logger)

	repos := Repos{
		releasesRepo: releasesRepo,

		releaseVersionsRepo: bhrelsrepo.NewConcreteReleaseVersionsRepository(options.ReleasesIndexDir, fs, logger),
		jobsRepo:            bhjobsrepo.NewConcreteJobsRepository(relVerFactory, logger),

		s3StemcellsRepo:   bhstemsrepo.NewS3StemcellsRepository(options.StemcellsIndexDir, stemcellNotesRepo, fs, logger),
		stemcellNotesRepo: stemcellNotesRepo,
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
