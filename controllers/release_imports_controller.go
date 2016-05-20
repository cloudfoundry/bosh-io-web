package controllers

import (
	"net/http"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	mart "github.com/go-martini/martini"
	martrend "github.com/martini-contrib/render"

	bhimpsrepo "github.com/cppforlife/bosh-hub/release/importsrepo"
)

type ReleaseImportsController struct {
	importsRepo  bhimpsrepo.ImportsRepository
	privateToken string

	indexTmpl string
	errorTmpl string

	logger boshlog.Logger
}

func NewReleaseImportsController(
	importsRepo bhimpsrepo.ImportsRepository,
	privateToken string,
	logger boshlog.Logger,
) ReleaseImportsController {
	return ReleaseImportsController{
		importsRepo:  importsRepo,
		privateToken: privateToken,

		indexTmpl: "release_imports/index",
		errorTmpl: "error",

		logger: logger,
	}
}

type releaseImportsControllerIndexPage struct {
	Imports      []bhimpsrepo.ImportRec
	PrivateToken string
}

func (c ReleaseImportsController) Index(r martrend.Render) {
	importRecs, err := c.importsRepo.ListAll()
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	page := releaseImportsControllerIndexPage{
		Imports:      importRecs,
		PrivateToken: c.privateToken,
	}

	r.HTML(200, c.indexTmpl, page)
}

func (c ReleaseImportsController) Delete(req *http.Request, r martrend.Render, params mart.Params) {
	importRec := bhimpsrepo.ImportRec{
		RelSource: req.FormValue("relSource"),
		Version:   req.FormValue("version"),
	}

	err := c.importsRepo.Remove(importRec)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	r.Redirect(c.privateToken + "/release_imports")
}
