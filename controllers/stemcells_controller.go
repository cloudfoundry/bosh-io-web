package controllers

import (
	"net/http"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	martrend "github.com/martini-contrib/render"

	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
	bhstemui "github.com/cppforlife/bosh-hub/ui/stemcell"
)

var (
	UniqueVersionStemcellsLimit int = 30
)

type StemcellsController struct {
	stemcellsRepo bhstemsrepo.StemcellsRepository

	indexTmpl string
	errorTmpl string

	logTag string
	logger boshlog.Logger
}

func NewStemcellsController(
	stemcellsRepo bhstemsrepo.StemcellsRepository,
	logger boshlog.Logger,
) StemcellsController {
	return StemcellsController{
		stemcellsRepo: stemcellsRepo,

		indexTmpl: "stemcells/index",
		errorTmpl: "error",

		logTag: "StemcellsController",
		logger: logger,
	}
}

type stemcellsControllerIndexPage struct {
	UniqueVersionStemcells []*bhstemui.SameVersionStemcells
}

func (c StemcellsController) Index(req *http.Request, r martrend.Render) {
	stemcells, err := c.stemcellsRepo.FindAll()
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	var limit *int

	// Limit number of stemcells shown by default
	if len(req.URL.Query().Get("all")) == 0 {
		limit = &UniqueVersionStemcellsLimit
	}

	page := stemcellsControllerIndexPage{
		UniqueVersionStemcells: bhstemui.NewUniqueVersionStemcells(stemcells, limit),
	}

	r.HTML(200, c.indexTmpl, page)
}
