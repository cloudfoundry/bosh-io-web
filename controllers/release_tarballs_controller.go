package controllers

import (
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	mart "github.com/go-martini/martini"
	martrend "github.com/martini-contrib/render"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

type ReleaseTarballsController struct {
	releasesRepo bhrelsrepo.ReleasesRepository

	errorTmpl string

	logger boshlog.Logger
}

func NewReleaseTarballsController(
	releasesRepo bhrelsrepo.ReleasesRepository,
	logger boshlog.Logger,
) ReleaseTarballsController {
	return ReleaseTarballsController{
		releasesRepo: releasesRepo,

		errorTmpl: "error",

		logger: logger,
	}
}

// Show uses '_1' param as release source and 'v' param as release version
func (c ReleaseTarballsController) Download(req *http.Request, r martrend.Render, params mart.Params) {
	relSource := params["_1"]

	if len(relSource) == 0 {
		err := bosherr.Error("Param 'source' must be non-empty")
		r.HTML(400, c.errorTmpl, err)
		return
	}

	relVersion := req.URL.Query().Get("v")

	var relVerRec bhrelsrepo.ReleaseVersionRec
	var err error

	if relVersion == "" {
		relVerRec, err = c.releasesRepo.FindLatest(relSource)
		if err != nil {
			r.HTML(500, c.errorTmpl, err)
			return
		}
	} else {
		relVerRec, err = c.releasesRepo.Find(relSource, relVersion)
		if err != nil {
			r.HTML(500, c.errorTmpl, err)
			return
		}
	}

	relTarRec, err := relVerRec.Tarball()
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	url, err := relTarRec.ActualDownloadURL()
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	r.Redirect(url)
}
