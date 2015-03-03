package controllers

import (
	"net/http"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	martrend "github.com/martini-contrib/render"

	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
	bhstemui "github.com/cppforlife/bosh-hub/ui/stemcell"
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
	DistroGroups bhstemui.DistroGroups
	Filter       bhstemui.StemcellFilter
}

func (c StemcellsController) Index(req *http.Request, r martrend.Render) {
	filter := bhstemui.StemcellFilter{Name: req.URL.Query().Get("name")}

	stemcells, err := c.stemcellsRepo.FindAll(filter.Name)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	// Show either groups of stemcells by OS or for a specific stemcell name
	distroGroups := bhstemui.NewDistroGroups(stemcells, filter)

	r.HTML(200, c.indexTmpl, stemcellsControllerIndexPage{distroGroups, filter})
}

func (c StemcellsController) APIV1Index(req *http.Request, r martrend.Render) {
	filter := bhstemui.StemcellFilter{Name: req.URL.Query().Get("name")}

	if len(filter.Name) == 0 {
		r.JSON(400, map[string]string{"error": "Param 'name' must be non-empty"})
		return
	}

	stemcells, err := c.stemcellsRepo.FindAll(filter.Name)
	if err != nil {
		r.JSON(500, map[string]string{"error": err.Error()})
		return
	}

	// Show list of latest versions for the specific stemcell name
	uniqVerStems := bhstemui.NewUniqueVersionStemcells(stemcells, filter)

	r.JSON(200, uniqVerStems.ForAPI())
}
