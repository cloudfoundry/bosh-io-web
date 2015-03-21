package controllers

import (
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
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
		err := bosherr.New("Param 'source' must be non-empty")
		r.HTML(500, c.errorTmpl, err)
		return
	}

	relVersion := req.URL.Query().Get("v")

	var relVerRec bhrelsrepo.ReleaseVersionRec
	var err error

	if relVersion == "" {
		var found bool

		relVerRec, found, err = c.releasesRepo.FindLatest(relSource)
		if err != nil {
			r.HTML(500, c.errorTmpl, err)
			return
		}

		if !found {
			err := bosherr.New("Latest release is not found")
			r.HTML(404, c.errorTmpl, err)
			return
		}
	} else {
		relVerRec, err = c.releasesRepo.Find(relSource, relVersion)
		if err != nil {
			r.HTML(500, c.errorTmpl, err)
			return
		}
	}

	relTarRec, found, err := relVerRec.Tarball()
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	if !found {
		err := bosherr.New("Release tarball is not found")
		r.HTML(404, c.errorTmpl, err)
		return
	}

	url, err := relTarRec.ActualDownloadURL()
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	r.Redirect(url)
}
