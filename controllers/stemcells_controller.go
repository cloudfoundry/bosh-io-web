package controllers

import (
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	semiver "github.com/cppforlife/go-semi-semantic/version"
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
	filter := bhstemui.StemcellFilter{
		Name:                     params["_1"],
		IncludeAll:               c.isParamTrue(req.URL.Query().Get("all")),
		IncludeDeprecatedDistros: true,
	}

	stemcells, err := c.stemcellsRepo.FindAll(filter.Name)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	// Show either groups of stemcells by OS or for a specific stemcell name
	distroGroups := bhstemui.NewDistroGroups(stemcells, filter)

	r.HTML(200, c.indexTmpl, stemcellsControllerIndexPage{distroGroups, filter})
}

// Show uses '_1' param as stemcell name and 'v' param as release version
func (c StemcellsController) Download(req *http.Request, r martrend.Render, params mart.Params) {
	relSource := params["_1"]

	if len(relSource) == 0 {
		err := bosherr.New("Param 'source' must be non-empty")
		r.HTML(400, c.errorTmpl, err)
		return
	}

	stemcells, err := c.stemcellsRepo.FindAll(relSource)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	filter := bhstemui.StemcellFilter{Name: relSource, IncludeAll: true}

	// List of versions for the specific stemcell name
	// todo really filtering below should be part of the repo instead of ui
	uniqVerStems := bhstemui.NewUniqueVersionStemcells(stemcells, filter)

	sortedStemcells := uniqVerStems.ForAPI()

	relVersion := req.URL.Query().Get("v")

	preferLight := true

	lightParam := req.URL.Query().Get("light")
	if lightParam != "" {
		preferLight = c.isParamTrue(lightParam)
	}

	// todo remove china stemcell once consolidated into light stemcells
	chinaParam := req.URL.Query().Get("china")
	mustBeForChina := c.isParamTrue(chinaParam)

	if relVersion == "" {
		if len(sortedStemcells) == 0 {
			err := bosherr.New("Latest stemcell is not found")
			r.HTML(404, c.errorTmpl, err)
			return
		}

		url, err := sortedStemcells[0].ActualDownloadURL(preferLight, mustBeForChina)
		if err != nil {
			r.HTML(404, c.errorTmpl, err)
			return
		}

		r.Redirect(url)
		return
	}

	ver, err := semiver.NewVersionFromString(relVersion)
	if err != nil {
		err = bosherr.New("Version '%s' is not valid: %s", relVersion, err)
		r.HTML(400, c.errorTmpl, err)
		return
	}

	for _, stemcell := range sortedStemcells {
		if stemcell.Version.IsEq(ver) {
			url, err := stemcell.ActualDownloadURL(preferLight, mustBeForChina)
			if err != nil {
				r.HTML(404, c.errorTmpl, err)
				return
			}

			r.Redirect(url)
			return
		}
	}

	err = bosherr.New("Stemcell version '%s' is not found", relVersion)
	r.HTML(404, c.errorTmpl, err)
}

func (c StemcellsController) isParamTrue(s string) bool {
	return s == "1" || s == "true" || s == "t"
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
