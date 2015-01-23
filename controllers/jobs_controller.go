package controllers

import (
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	mart "github.com/go-martini/martini"
	martrend "github.com/martini-contrib/render"

	bhjobsrepo "github.com/cppforlife/bosh-hub/release/jobsrepo"
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhjobui "github.com/cppforlife/bosh-hub/ui/job"
	bhrelui "github.com/cppforlife/bosh-hub/ui/release"
)

type JobsController struct {
	releaseVersionsRepo bhrelsrepo.ReleaseVersionsRepository
	jobsRepo            bhjobsrepo.JobsRepository

	showTmpl  string
	errorTmpl string

	logTag string
	logger boshlog.Logger
}

func NewJobsController(
	releaseVersionsRepo bhrelsrepo.ReleaseVersionsRepository,
	jobsRepo bhjobsrepo.JobsRepository,
	logger boshlog.Logger,
) JobsController {
	return JobsController{
		releaseVersionsRepo: releaseVersionsRepo,
		jobsRepo:            jobsRepo,

		showTmpl:  "jobs/show",
		errorTmpl: "error",

		logTag: "JobsController",
		logger: logger,
	}
}

func (c JobsController) Show(req *http.Request, r martrend.Render, params mart.Params) {
	relSource, relVersion, jobName, err := c.extractShowParams(req, params)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	c.logger.Debug(c.logTag, "Release source '%s'", relSource)

	relVerRec := bhrelsrepo.ReleaseVersionRec{
		Source:     relSource,
		VersionRaw: relVersion,
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

	if !found {
		err := bosherr.New("Release jobs are not found")
		r.HTML(404, c.errorTmpl, err)
		return
	}

	viewRel := bhrelui.NewRelease(relSource, rel)

	for _, relJob := range relJobs {
		if relJob.Name == jobName {
			r.HTML(200, c.showTmpl, bhjobui.NewJob(relJob, viewRel))
			return
		}
	}

	err = bosherr.New("Release job '%s' is not found", jobName)
	r.HTML(404, c.errorTmpl, err)
}

func (c JobsController) extractShowParams(req *http.Request, params mart.Params) (string, string, string, error) {
	relSource := req.URL.Query().Get("source")

	if len(relSource) == 0 {
		return "", "", "", bosherr.New("Param 'source' must be non-empty")
	}

	relVersion := req.URL.Query().Get("version")

	if len(relVersion) == 0 {
		return "", "", "", bosherr.New("Param 'version' must be non-empty")
	}

	jobName := params["name"]

	if len(jobName) == 0 {
		return "", "", "", bosherr.New("Param 'name' must be non-empty")
	}

	return relSource, relVersion, jobName, nil
}
