package controllers

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	mart "github.com/go-martini/martini"
	martrend "github.com/martini-contrib/render"

	bhbibrepo "github.com/cppforlife/bosh-hub/bosh-init-bin/repo"
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhrelui "github.com/cppforlife/bosh-hub/ui/release"
)

var (
	// Be very conservative about which characters are allowed from the start to the end
	docsControllerPageRegexp = regexp.MustCompile(`\A[a-zA-Z0-9\-/]+\z`)
	docsControllerGithubURL  = "https://github.com/cloudfoundry/docs-bosh/blob/master"
)

type DocsController struct {
	releasesRepo        bhrelsrepo.ReleasesRepository
	releaseVersionsRepo bhrelsrepo.ReleaseVersionsRepository
	boshInitBinsRepo    bhbibrepo.Repository

	defaultTmpl string
	errorTmpl   string

	logTag string
	logger boshlog.Logger
}

func NewDocsController(
	releasesRepo bhrelsrepo.ReleasesRepository,
	releaseVersionsRepo bhrelsrepo.ReleaseVersionsRepository,
	boshInitBinsRepo bhbibrepo.Repository,
	logger boshlog.Logger,
) DocsController {
	return DocsController{
		releasesRepo:        releasesRepo,
		releaseVersionsRepo: releaseVersionsRepo,
		boshInitBinsRepo:    boshInitBinsRepo,

		defaultTmpl: "index",
		errorTmpl:   "error",

		logTag: "DocsController",
		logger: logger,
	}
}

type docPage struct {
	ContributeChangesURL string
}

type installBoshInitPage struct {
	docPage
	LatestBoshInitBinGroups []bhbibrepo.BinaryGroup
}

type initManifestPage struct {
	docPage
	Releases []bhrelui.Release
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

	if tmpl == "install-bosh-init" {
		binGroups, err := c.boshInitBinsRepo.FindLatest()
		if err != nil {
			return nil, err
		}

		return installBoshInitPage{docPage: page, LatestBoshInitBinGroups: binGroups}, nil
	}

	// init-<cpi> pages have an example manifest which lists latest releases
	cpi, found := bhrelui.KnownCPIs.FindByDocPage(tmpl)
	if found {
		return newInitManifestPage(page, cpi, c.releasesRepo, c.releaseVersionsRepo)
	}

	return page, nil
}

func newInitManifestPage(
	docPage docPage,
	cpiRef bhrelui.ReleaseRef,
	releasesRepo bhrelsrepo.ReleasesRepository,
	releaseVersionsRepo bhrelsrepo.ReleaseVersionsRepository,
) (initManifestPage, error) {
	page := initManifestPage{docPage: docPage}

	sources := []bhrelui.Source{bhrelui.BOSH.Source, cpiRef.Source}

	for _, source := range sources {
		relVerRec, found, err := releasesRepo.FindLatest(source.Full())
		if err != nil {
			return page, err
		} else if !found {
			return page, bosherr.New("Latest release '%s' is not found", source.Full())
		}

		rel, found, err := releaseVersionsRepo.Find(relVerRec)
		if err != nil {
			return page, err
		} else if !found {
			return page, bosherr.New("Release '%s' is not found", source.Full())
		}

		page.Releases = append(page.Releases, bhrelui.NewRelease(relVerRec, rel))
	}

	return page, nil
}
