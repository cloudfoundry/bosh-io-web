package controllers

import (
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	mart "github.com/go-martini/martini"
	martrend "github.com/martini-contrib/render"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhrelui "github.com/cppforlife/bosh-hub/ui/release"
)

type PackagesController struct {
	releaseVersionsRepo bhrelsrepo.ReleaseVersionsRepository

	showTmpl  string
	errorTmpl string

	runner boshsys.CmdRunner

	logTag string
	logger boshlog.Logger
}

func NewPackagesController(
	releaseVersionsRepo bhrelsrepo.ReleaseVersionsRepository,
	runner boshsys.CmdRunner,
	logger boshlog.Logger,
) PackagesController {
	return PackagesController{
		releaseVersionsRepo: releaseVersionsRepo,

		showTmpl:  "packages/show",
		errorTmpl: "error",

		runner: runner,

		logTag: "PackagesController",
		logger: logger,
	}
}

func (c PackagesController) Show(req *http.Request, r martrend.Render, params mart.Params) {
	relSource, relVersion, pkgName, err := c.extractShowParams(req, params)
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

	viewRel := bhrelui.NewRelease(relSource, rel)

	for _, viewPkg := range viewRel.Packages {
		if viewPkg.Name == pkgName {
			r.HTML(200, c.showTmpl, viewPkg)
			return
		}
	}

	err = bosherr.New("Release package '%s' is not found", pkgName)
	r.HTML(404, c.errorTmpl, err)
}

func (c PackagesController) extractShowParams(req *http.Request, params mart.Params) (string, string, string, error) {
	relSource := req.URL.Query().Get("source")

	if len(relSource) == 0 {
		return "", "", "", bosherr.New("Param 'source' must be non-empty")
	}

	relVersion := req.URL.Query().Get("version")

	if len(relVersion) == 0 {
		return "", "", "", bosherr.New("Param 'version' must be non-empty")
	}

	pkgName := params["name"]

	if len(pkgName) == 0 {
		return "", "", "", bosherr.New("Param 'name' must be non-empty")
	}

	return relSource, relVersion, pkgName, nil
}
