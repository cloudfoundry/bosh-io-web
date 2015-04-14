package controllers

import (
	"fmt"
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	mart "github.com/go-martini/martini"
	martrend "github.com/martini-contrib/render"

	bhjobsrepo "github.com/cppforlife/bosh-hub/release/jobsrepo"
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
	bhjobui "github.com/cppforlife/bosh-hub/ui/job"
	bhmiscui "github.com/cppforlife/bosh-hub/ui/misc"
	bhrelui "github.com/cppforlife/bosh-hub/ui/release"
	bhstemui "github.com/cppforlife/bosh-hub/ui/stemcell"
)

type ReleasesController struct {
	releasesRepo        bhrelsrepo.ReleasesRepository
	releaseVersionsRepo bhrelsrepo.ReleaseVersionsRepository
	jobsRepo            bhjobsrepo.JobsRepository

	stemcellsRepo bhstemsrepo.StemcellsRepository

	indexTmpl        string
	showVersionsTmpl string
	showVersionTmpl  string
	graphTmpl        string
	errorTmpl        string

	runner boshsys.CmdRunner

	logTag string
	logger boshlog.Logger
}

func NewReleasesController(
	releasesRepo bhrelsrepo.ReleasesRepository,
	releaseVersionsRepo bhrelsrepo.ReleaseVersionsRepository,
	jobsRepo bhjobsrepo.JobsRepository,
	stemcellsRepo bhstemsrepo.StemcellsRepository,
	runner boshsys.CmdRunner,
	logger boshlog.Logger,
) ReleasesController {
	return ReleasesController{
		releasesRepo:        releasesRepo,
		releaseVersionsRepo: releaseVersionsRepo,
		jobsRepo:            jobsRepo,

		stemcellsRepo: stemcellsRepo,

		indexTmpl:        "releases/index",
		showVersionsTmpl: "releases/show_versions",
		showVersionTmpl:  "releases/show_version",
		graphTmpl:        "releases/graph",
		errorTmpl:        "error",

		runner: runner,

		logTag: "ReleasesController",
		logger: logger,
	}
}

type releasesControllerHomePage struct {
	UniqueSourceReleases   bhrelui.UniqueSourceReleases
	LatestVersionStemcells *bhstemui.SameVersionStemcells
}

type releasesControllerIndexPage struct {
	UniqueSources bhrelui.UniqueSources
}

func (c ReleasesController) Index(r martrend.Render) {
	sources, err := c.releasesRepo.ListAll()
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	page := releasesControllerIndexPage{
		UniqueSources: bhrelui.NewUniqueSources(sources),
	}

	r.HTML(200, c.indexTmpl, page)
}

// Show uses '_1' param as release source and 'version' param as release version
func (c ReleasesController) Show(req *http.Request, r martrend.Render, params mart.Params) {
	relSource := params["_1"]

	if len(relSource) == 0 {
		err := bosherr.New("Param 'source' must be non-empty")
		r.HTML(500, c.errorTmpl, err)
		return
	}

	c.logger.Debug(c.logTag, "Release source '%s'", relSource)

	relVersion := req.URL.Query().Get("version")

	graph := req.URL.Query().Get("graph")

	if relVersion == "" {
		c.showMultipleReleases(r, relSource)
	} else {
		var tmpl string

		if len(graph) > 0 {
			tmpl = c.graphTmpl
		} else {
			tmpl = c.showVersionTmpl
		}

		c.showSingleRelease(r, relSource, relVersion, tmpl)
	}
}

func (c ReleasesController) showMultipleReleases(r martrend.Render, relSource string) {
	relVerRecs, found, err := c.releasesRepo.FindAll(relSource)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	if !found {
		err := bosherr.New("Release '%s' is not found", relSource)
		r.HTML(404, c.errorTmpl, err)
		return
	}

	viewRels := bhrelui.NewSameSourceReleases(bhrelsrepo.Source{Full: relSource}, relVerRecs)

	r.HTML(200, c.showVersionsTmpl, viewRels)
}

func (c ReleasesController) showSingleRelease(r martrend.Render, relSource, relVersion, tmpl string) {
	relVerRec, err := c.releasesRepo.Find(relSource, relVersion)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	rel, found, err := c.releaseVersionsRepo.Find(relVerRec)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	if !found {
		err := bosherr.New("Release '%s' is not found", relSource)
		r.HTML(404, c.errorTmpl, err)
		return
	}

	relJobs, found, err := c.jobsRepo.FindAll(relVerRec)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	viewRel := bhrelui.NewRelease(relVerRec, rel)

	viewJobs := []bhjobui.Job{}

	for _, relJob := range relJobs {
		viewJobs = append(viewJobs, bhjobui.NewJob(relJob, viewRel))
	}

	viewRel.Graph = bhmiscui.NewReleaseGraph(viewRel.Packages, viewJobs, c.runner, c.logger)

	r.HTML(200, tmpl, &viewRel)
}

func (c ReleasesController) APIV1Index(req *http.Request, r martrend.Render, params mart.Params) {
	relSource := params["_1"]

	if len(relSource) == 0 {
		r.JSON(400, map[string]string{"error": "Param 'source' must be non-empty"})
		return
	}

	relVerRecs, found, err := c.releasesRepo.FindAll(relSource)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	if !found {
		r.JSON(404, map[string]string{"error": fmt.Sprintf("Release '%s' is not found", relSource)})
		return
	}

	// Show list of latest versions for the specific stemcell name
	viewRels := bhrelui.NewSameSourceReleases(bhrelsrepo.Source{Full: relSource}, relVerRecs)

	r.JSON(200, viewRels.ForAPI())
}
