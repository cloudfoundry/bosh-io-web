package controllers

import (
	"errors"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"

	bhbibrepo "github.com/cppforlife/bosh-hub/bosh-init-bin/repo"
	bhchecksrepo "github.com/cppforlife/bosh-hub/checksumsrepo"
	bhimperrsrepo "github.com/cppforlife/bosh-hub/release/importerrsrepo"
	bhimpsrepo "github.com/cppforlife/bosh-hub/release/importsrepo"
	bhjobsrepo "github.com/cppforlife/bosh-hub/release/jobsrepo"
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhwatchersrepo "github.com/cppforlife/bosh-hub/release/watchersrepo"
	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
)

type FactoryRepos interface {
	ReleasesRepo() bhrelsrepo.ReleasesRepository
	ReleaseVersionsRepo() bhrelsrepo.ReleaseVersionsRepository
	JobsRepo() bhjobsrepo.JobsRepository

	S3StemcellsRepo() bhstemsrepo.S3StemcellsRepository
	StemcellsRepo() bhstemsrepo.StemcellsRepository

	BoshInitBinsRepo() bhbibrepo.Repository

	ImportsRepo() bhimpsrepo.ImportsRepository
	ImportErrsRepo() bhimperrsrepo.ImportErrsRepository
	WatchersRepo() bhwatchersrepo.WatchersRepository

	ChecksumsRepo() bhchecksrepo.ChecksumsRepository
}

type Factory struct {
	HomeController HomeController
	DocsController DocsController

	ReleasesController        ReleasesController
	ReleaseTarballsController ReleaseTarballsController

	StemcellsController StemcellsController

	JobsController     JobsController
	PackagesController PackagesController

	ReleaseWatchersController   ReleaseWatchersController
	ReleaseImportsController    ReleaseImportsController
	ReleaseImportErrsController ReleaseImportErrsController

	ChecksumsController ChecksumsController

	privateURLPrefix string
}

func NewFactory(privateToken string, checksumPrivs []ChecksumReqMatch, r FactoryRepos, runner boshsys.CmdRunner, logger boshlog.Logger) (Factory, error) {
	privateToken = strings.TrimSpace(privateToken)

	if len(privateToken) < 10 {
		return Factory{}, errors.New("Expected private token to be at least 10 chars")
	} else {
		logger.Info("controllers.Factory", "Private token is '%s'", privateToken)
	}

	privateURLPrefix := "/" + privateToken

	factory := Factory{
		HomeController: NewHomeController(r.ReleasesRepo(), r.StemcellsRepo(), logger),

		DocsController: NewDocsController(
			r.ReleasesRepo(),
			r.ReleaseVersionsRepo(),
			r.BoshInitBinsRepo(),
			logger,
		),

		ReleasesController: NewReleasesController(
			r.ReleasesRepo(),
			r.ReleaseVersionsRepo(),
			r.JobsRepo(),
			r.StemcellsRepo(),
			runner,
			logger,
		),

		ReleaseTarballsController: NewReleaseTarballsController(r.ReleasesRepo(), logger),

		StemcellsController: NewStemcellsController(r.StemcellsRepo(), logger),

		JobsController:     NewJobsController(r.ReleasesRepo(), r.ReleaseVersionsRepo(), r.JobsRepo(), logger),
		PackagesController: NewPackagesController(r.ReleasesRepo(), r.ReleaseVersionsRepo(), runner, logger),

		ReleaseWatchersController:   NewReleaseWatchersController(r.WatchersRepo(), privateURLPrefix, logger),
		ReleaseImportsController:    NewReleaseImportsController(r.ImportsRepo(), privateURLPrefix, logger),
		ReleaseImportErrsController: NewReleaseImportErrsController(r.ImportErrsRepo(), privateURLPrefix, logger),

		ChecksumsController: NewChecksumsController(checksumPrivs, r.ChecksumsRepo(), logger),

		privateURLPrefix: privateURLPrefix,
	}

	return factory, nil
}

func (f Factory) PrivateURL(ending string) string {
	return f.privateURLPrefix + ending
}
