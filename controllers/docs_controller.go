package controllers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	mart "github.com/go-martini/martini"
	martrend "github.com/martini-contrib/render"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
	bhrelui "github.com/cppforlife/bosh-hub/ui/release"
	bhstemui "github.com/cppforlife/bosh-hub/ui/stemcell"
)

var (
	// Be very conservative about which characters are allowed from the start to the end
	docsControllerPageRegexp = regexp.MustCompile(`\A[a-zA-Z0-9\-/]+\z`)
	docsControllerGithubURL  = "https://github.com/cloudfoundry/docs-bosh/blob/master"
)

type DocsController struct {
	releasesRepo        bhrelsrepo.ReleasesRepository
	releaseVersionsRepo bhrelsrepo.ReleaseVersionsRepository
	stemcellsRepo       bhstemsrepo.StemcellsRepository

	defaultTmpl string
	errorTmpl   string

	logTag string
	logger boshlog.Logger
}

func NewDocsController(
	releasesRepo bhrelsrepo.ReleasesRepository,
	releaseVersionsRepo bhrelsrepo.ReleaseVersionsRepository,
	stemcellsRepo bhstemsrepo.StemcellsRepository,
	logger boshlog.Logger,
) DocsController {
	return DocsController{
		releasesRepo:        releasesRepo,
		releaseVersionsRepo: releaseVersionsRepo,
		stemcellsRepo:       stemcellsRepo,

		defaultTmpl: "index",
		errorTmpl:   "error",

		logTag: "DocsController",
		logger: logger,
	}
}

type docPage struct {
	ContributeChangesURL string
}

type initManifestPage struct {
	docPage
	Releases []bhrelui.Release
	Stemcell bhstemui.Stemcell
}

func (c DocsController) Page(r martrend.Render, params mart.Params) {
	tmpl := strings.TrimSuffix(params["_1"], ".html")

	if len(tmpl) == 0 {
		tmpl = c.defaultTmpl
	}

	if !docsControllerPageRegexp.MatchString(tmpl) {
		err := errors.New("Invalid page")
		r.HTML(500, c.errorTmpl, err)
		return
	}

	page, err := c.findPage(tmpl)
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	r.HTML(200, "docs/"+tmpl, page)
}

func (c DocsController) findPage(tmpl string) (interface{}, error) {
	page := docPage{
		ContributeChangesURL: fmt.Sprintf("%s/%s.html.md.erb", docsControllerGithubURL, tmpl),
	}

	// init-<cpi> pages have an example manifest which lists latest releases
	cpiRef, found := bhrelui.KnownCPIs.FindByDocPage(tmpl)
	if found {
		// populate example manifest with latest stemcell
		stemcellRef, found := bhstemui.KnownStemcells.FindByDocPage(tmpl)
		if found {
			return c.newInitManifestPage(page, cpiRef, stemcellRef)
		}
	}

	return page, nil
}

func (c DocsController) newInitManifestPage(docPage docPage, cpiRef bhrelui.ReleaseRef, stemcellRef bhstemui.StemcellRef) (initManifestPage, error) {
	page := initManifestPage{docPage: docPage}

	sources := []bhrelui.Source{bhrelui.BOSH.Source, cpiRef.Source}

	for _, source := range sources {
		relVerRec, err := c.releasesRepo.FindLatest(source.Full())
		if err != nil {
			return page, err
		}

		rel, err := c.releaseVersionsRepo.Find(relVerRec)
		if err != nil {
			return page, err
		}

		page.Releases = append(page.Releases, bhrelui.NewRelease(relVerRec, rel))
	}

	stemcells, err := c.stemcellsRepo.FindAll(stemcellRef.ManifestName)
	if err != nil {
		return page, err
	}

	uniqVerStems := bhstemui.NewUniqueVersionStemcells(stemcells, bhstemui.StemcellFilter{})

	page.Stemcell = uniqVerStems.ForAPI()[0]

	return page, nil
}
