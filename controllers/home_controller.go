package controllers

import (
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	martrend "github.com/martini-contrib/render"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
	bhrelui "github.com/cppforlife/bosh-hub/ui/release"
	bhstemui "github.com/cppforlife/bosh-hub/ui/stemcell"
)

type HomeController struct {
	releasesRepo  bhrelsrepo.ReleasesRepository
	stemcellsRepo bhstemsrepo.StemcellsRepository

	homeTmpl  string
	errorTmpl string

	logTag string
	logger boshlog.Logger
}

func NewHomeController(
	releasesRepo bhrelsrepo.ReleasesRepository,
	stemcellsRepo bhstemsrepo.StemcellsRepository,
	logger boshlog.Logger,
) HomeController {
	return HomeController{
		releasesRepo:  releasesRepo,
		stemcellsRepo: stemcellsRepo,

		homeTmpl:  "home/home",
		errorTmpl: "error",

		logTag: "HomeController",
		logger: logger,
	}
}

type homeControllerPage struct {
	UniqueSourceReleases bhrelui.UniqueSourceReleases
	StemcellDistroGroups bhstemui.DistroGroups
}

func (c HomeController) Home(r martrend.Render) {
	relVerRecs, err := c.releasesRepo.ListCurated()
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	stemcells, err := c.stemcellsRepo.FindAll("")
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	page := homeControllerPage{
		UniqueSourceReleases: bhrelui.NewUniqueSourceReleases(relVerRecs),

		// Show either groups of stemcells by OS
		StemcellDistroGroups: bhstemui.NewDistroGroups(stemcells, bhstemui.StemcellFilter{}),
	}

	r.HTML(200, c.homeTmpl, page)
}
