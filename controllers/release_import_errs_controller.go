package controllers

import (
	"net/http"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	mart "github.com/go-martini/martini"
	martrend "github.com/martini-contrib/render"

	bhimperrsrepo "github.com/cppforlife/bosh-hub/release/importerrsrepo"
	bhimpsrepo "github.com/cppforlife/bosh-hub/release/importsrepo"
)

type ReleaseImportErrsController struct {
	importErrsRepo bhimperrsrepo.ImportErrsRepository
	privateToken   string

	indexTmpl string
	errorTmpl string

	logger boshlog.Logger
}

func NewReleaseImportErrsController(
	importErrsRepo bhimperrsrepo.ImportErrsRepository,
	privateToken string,
	logger boshlog.Logger,
) ReleaseImportErrsController {
	return ReleaseImportErrsController{
		importErrsRepo: importErrsRepo,
		privateToken:   privateToken,

		indexTmpl: "release_import_errs/index",
		errorTmpl: "error",

		logger: logger,
	}
}

type releaseImportErrsControllerIndexPage struct {
	ImportErrs   []bhimperrsrepo.ImportErrRec
	PrivateToken string
}

func (c ReleaseImportErrsController) Index(r martrend.Render) {
	importErrRecs, err := c.importErrsRepo.ListAll()
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	page := releaseImportErrsControllerIndexPage{
		ImportErrs:   importErrRecs,
		PrivateToken: c.privateToken,
	}

	r.HTML(200, c.indexTmpl, page)
}

func (c ReleaseImportErrsController) Delete(req *http.Request, r martrend.Render, params mart.Params) {
	importRec := bhimpsrepo.ImportRec{
		RelSource: req.FormValue("relSource"),
		Version:   req.FormValue("version"),
	}

	err := c.importErrsRepo.Remove(importRec)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	r.Redirect(c.privateToken + "/release_import_errs")
}
