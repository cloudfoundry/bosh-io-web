package controllers

import (
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	mart "github.com/go-martini/martini"
	martrend "github.com/martini-contrib/render"

	bhjobsrepo "github.com/bosh-io/web/release/jobsrepo"
	bhrelsrepo "github.com/bosh-io/web/release/releasesrepo"
	bhstemsrepo "github.com/bosh-io/web/stemcell/stemsrepo"
	bhjobui "github.com/bosh-io/web/ui/job"
	bhmiscui "github.com/bosh-io/web/ui/misc"
	bhrelui "github.com/bosh-io/web/ui/release"
	bhstemui "github.com/bosh-io/web/ui/stemcell"
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
	relSource, showAll, relVersion, showGraph, err := c.extractShowParams(req, params)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	c.logger.Debug(c.logTag, "Release source '%s'", relSource)

	if len(showAll) > 0 {
		c.showMultipleReleases(r, relSource)
	} else {
		var tmpl string

		if len(showGraph) > 0 {
			tmpl = c.graphTmpl
		} else {
			tmpl = c.showVersionTmpl
		}

		c.showSingleRelease(r, relSource, relVersion, tmpl)
	}
}

func (c ReleasesController) extractShowParams(req *http.Request, params mart.Params) (string, string, string, string, error) {
	relSource := params["_1"]

	if len(relSource) == 0 {
		return "", "", "", "", bosherr.Error("Param 'source' must be non-empty")
	}

	showAll := req.URL.Query().Get("all")

	relVersion := req.URL.Query().Get("version")

	showGraph := req.URL.Query().Get("graph")

	return relSource, showAll, relVersion, showGraph, nil
}

func (c ReleasesController) showMultipleReleases(r martrend.Render, relSource string) {
	relVerRecs, err := c.releasesRepo.FindAll(relSource)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	var relName string

	// Fetch full release details for one of the versions to get real release name
	if len(relVerRecs) > 0 {
		rel, err := c.releaseVersionsRepo.Find(relVerRecs[0])
		if err != nil {
			r.HTML(500, c.errorTmpl, err)
			return
		}

		relName = rel.Name
	}

	viewRels := bhrelui.NewSameSourceReleases(bhrelsrepo.Source{Full: relSource}, relVerRecs, relName)

	r.HTML(200, c.showVersionsTmpl, viewRels)
}

func (c ReleasesController) showSingleRelease(r martrend.Render, relSource, relVersion, tmpl string) {
	var err error

	var relVerRec bhrelsrepo.ReleaseVersionRec

	if len(relVersion) > 0 {
		relVerRec, err = c.releasesRepo.Find(relSource, relVersion)
		if err != nil {
			r.HTML(500, c.errorTmpl, err)
			return
		}
	} else {
		relVerRec, err = c.releasesRepo.FindLatest(relSource)
		if err != nil {
			r.HTML(500, c.errorTmpl, err)
			return
		}
	}

	rel, err := c.releaseVersionsRepo.Find(relVerRec)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	relJobs, err := c.jobsRepo.FindAll(relVerRec)
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

	relVerRecs, err := c.releasesRepo.FindAll(relSource)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	// Show list of latest versions for the specific stemcell name
	viewRels := bhrelui.NewSameSourceReleases(bhrelsrepo.Source{Full: relSource}, relVerRecs, "")

	r.JSON(200, viewRels.ForAPI())
}
