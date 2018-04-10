package controllers

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bhjobsrepo "github.com/bosh-io/web/release/jobsrepo"
	bhrelsrepo "github.com/bosh-io/web/release/releasesrepo"
	bhstemsrepo "github.com/bosh-io/web/stemcell/stemsrepo"
)

type FactoryRepos interface {
	ReleasesRepo() bhrelsrepo.ReleasesRepository
	ReleaseVersionsRepo() bhrelsrepo.ReleaseVersionsRepository
	JobsRepo() bhjobsrepo.JobsRepository

	S3StemcellsRepo() bhstemsrepo.S3StemcellsRepository
	StemcellsRepo() bhstemsrepo.StemcellsRepository
}

type Factory struct {
	RedirectsController RedirectsController

	DocsController DocsController

	ReleasesController        ReleasesController
	ReleaseTarballsController ReleaseTarballsController

	StemcellsController StemcellsController

	JobsController     JobsController
	PackagesController PackagesController
}

func NewFactory(redirects RedirectsConfig, r FactoryRepos, runner boshsys.CmdRunner, logger boshlog.Logger) (Factory, error) {
	factory := Factory{
		RedirectsController: NewRedirectsController(redirects),

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
	}

	return factory, nil
}
