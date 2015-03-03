package controllers

import (
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	mart "github.com/go-martini/martini"
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

func (c StemcellsController) Index(req *http.Request, r martrend.Render, params mart.Params) {
	filter := bhstemui.StemcellFilter{Name: params["_1"]}

	stemcells, err := c.stemcellsRepo.FindAll(filter.Name)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	// Show either groups of stemcells by OS or for a specific stemcell name
	distroGroups := bhstemui.NewDistroGroups(stemcells, filter)

	r.HTML(200, c.indexTmpl, stemcellsControllerIndexPage{distroGroups, filter})
}

// Show uses '_1' param as stemcell name
func (c StemcellsController) Download(req *http.Request, r martrend.Render, params mart.Params) {
	filter := bhstemui.StemcellFilter{Name: params["_1"]}

	if len(filter.Name) == 0 {
		r.JSON(400, map[string]string{"error": "Stemcell name must be non-empty"})
		return
	}

	stemcells, err := c.stemcellsRepo.FindAll(filter.Name)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	// Show list of latest versions for the specific stemcell name
	uniqVerStems := bhstemui.NewUniqueVersionStemcells(stemcells, filter)

	sortedStemcells := uniqVerStems.ForAPI()

	if len(sortedStemcells) == 0 {
		err := bosherr.New("Latest stemcell is not found")
		r.HTML(404, c.errorTmpl, err)
		return
	}

	r.Redirect(sortedStemcells[0].ActualDownloadURL())
}

func (c StemcellsController) APIV1Index(req *http.Request, r martrend.Render, params mart.Params) {
	filter := bhstemui.StemcellFilter{Name: params["_1"]}

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
