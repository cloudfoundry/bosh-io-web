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

	indexTmpl string
	errorTmpl string

	logger boshlog.Logger
}

func NewReleaseImportErrsController(
	importErrsRepo bhimperrsrepo.ImportErrsRepository,
	logger boshlog.Logger,
) ReleaseImportErrsController {
	return ReleaseImportErrsController{
		importErrsRepo: importErrsRepo,

		indexTmpl: "release_import_errs/index",
		errorTmpl: "error",

		logger: logger,
	}
}

type releaseImportErrsControllerIndexPage struct {
	ImportErrs []bhimperrsrepo.ImportErrRec
}

func (c ReleaseImportErrsController) Index(r martrend.Render) {
	importErrRecs, err := c.importErrsRepo.ListAll()
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	page := releaseImportErrsControllerIndexPage{
		ImportErrs: importErrRecs,
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

	r.Redirect("/release_import_errs")
}
