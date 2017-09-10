package controllers

import (
	"errors"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"

	bhchecksrepo "github.com/cppforlife/bosh-hub/checksumsrepo"
	bhjobsrepo "github.com/cppforlife/bosh-hub/release/jobsrepo"
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
)

type FactoryRepos interface {
	ReleasesRepo() bhrelsrepo.ReleasesRepository
	ReleaseVersionsRepo() bhrelsrepo.ReleaseVersionsRepository
	JobsRepo() bhjobsrepo.JobsRepository

	S3StemcellsRepo() bhstemsrepo.S3StemcellsRepository
	StemcellsRepo() bhstemsrepo.StemcellsRepository

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

	ChecksumsController ChecksumsController

	privateURLPrefix string
}

func NewFactory(privateToken string, checksumPrivs []ChecksumReqMatch, r FactoryRepos, runner boshsys.CmdRunner, logger boshlog.Logger) (Factory, error) {
	privateToken = strings.TrimSpace(privateToken)

	if len(privateToken) < 10 {
		return Factory{}, errors.New("Expected private token to be at least 10 chars")
	}

	privateURLPrefix := "/" + privateToken

	factory := Factory{
		HomeController: NewHomeController(r.ReleasesRepo(), r.StemcellsRepo(), logger),

		DocsController: NewDocsController(
			r.ReleasesRepo(),
			r.ReleaseVersionsRepo(),
			r.StemcellsRepo(),
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

		ChecksumsController: NewChecksumsController(checksumPrivs, r.ChecksumsRepo(), logger),

		privateURLPrefix: privateURLPrefix,
	}

	return factory, nil
}

func (f Factory) PrivateURL(ending string) string {
	return f.privateURLPrefix + ending
}
