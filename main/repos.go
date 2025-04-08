package main

import (
	"os"
	"os/signal"
	"syscall"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bhjobsrepo "github.com/bosh-io/web/release/jobsrepo"
	bhnotesrepo "github.com/bosh-io/web/release/notesrepo"
	bhrelsrepo "github.com/bosh-io/web/release/releasesrepo"
	bhreltarsrepo "github.com/bosh-io/web/release/releasetarsrepo"
	bhrelver "github.com/bosh-io/web/release/relver"
	bhs3 "github.com/bosh-io/web/s3"
	bhstemnotesrepo "github.com/bosh-io/web/stemcell/notesrepo"
	bhstemsrepo "github.com/bosh-io/web/stemcell/stemsrepo"
)

type ReposOptions struct {
	Type string // e.g. file, db

	Dir     string
	ConnURL string

	ReleasesDir      string
	ReleasesIndexDir string

	StemcellsIndexDirs      []string
	StemcellsLegacyIndexDir string

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
	cachingFS := NewCachingFileSystem(fs, logger)
	fs = cachingFS

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)

	go func() {
		for { //nolint:staticcheck
			select {
			case <-sigs:
				cachingFS.DropCache()
			}
		}
	}()

	var err error

	linkerOpts := options.ReleaseTarballLinker

	var linkerFactory bhs3.URLFactory

	switch { //nolint:staticcheck
	case linkerOpts.Type == "CloudFront":
		linkerFactory, err = bhs3.NewCDNURLFactory(linkerOpts.BaseURL, linkerOpts.KeyPairID, linkerOpts.PrivateKey)
	case linkerOpts.Type == "S3":
		linkerFactory = bhs3.NewDirectURLFactory(linkerOpts.BaseURL)
	default:
		err = bosherr.Errorf("Expected linker type '%s'", linkerOpts.Type)
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

	stemcellNotesRepo := bhstemnotesrepo.NewConcreteNotesRepository(
		options.StemcellsLegacyIndexDir, options.StemcellsIndexDirs, fs, logger)

	repos := Repos{
		releasesRepo: releasesRepo,

		releaseVersionsRepo: bhrelsrepo.NewConcreteReleaseVersionsRepository(relVerFactory, logger),
		jobsRepo:            bhjobsrepo.NewConcreteJobsRepository(relVerFactory, logger),

		s3StemcellsRepo: bhstemsrepo.NewS3StemcellsRepository(
			options.StemcellsLegacyIndexDir, options.StemcellsIndexDirs, stemcellNotesRepo, fs, logger),
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
