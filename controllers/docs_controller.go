package controllers

import (
	"errors"
	"regexp"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	mart "github.com/go-martini/martini"
	martrend "github.com/martini-contrib/render"

	bhbibrepo "github.com/cppforlife/bosh-hub/bosh-init-bin/repo"
)

var (
	// Be very conservative about which characters are allowed from the start to the end
	docsControllerPageRegexp = regexp.MustCompile(`\A[a-zA-Z0-9\-/]+\z`)
)

type DocsController struct {
	boshInitBinsRepo bhbibrepo.Repository

	defaultTmpl string
	errorTmpl   string

	logTag string
	logger boshlog.Logger
}

func NewDocsController(boshInitBinsRepo bhbibrepo.Repository, logger boshlog.Logger) DocsController {
	return DocsController{
		boshInitBinsRepo: boshInitBinsRepo,

		defaultTmpl: "index",
		errorTmpl:   "error",

		logTag: "DocsController",
		logger: logger,
	}
}

type installBoshInitPage struct {
	LatestBoshInitBinGroups []bhbibrepo.BinaryGroup
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
	var page interface{}

	if tmpl == "install-bosh-init" {
		binGroups, err := c.boshInitBinsRepo.FindLatest()
		if err != nil {
			return nil, err
		}

		page = installBoshInitPage{LatestBoshInitBinGroups: binGroups}
	}

	return page, nil
}
