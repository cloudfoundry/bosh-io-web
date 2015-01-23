package controllers

import (
	"errors"
	"regexp"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	mart "github.com/go-martini/martini"
	martrend "github.com/martini-contrib/render"
)

var (
	// Be very conservative about which characters are allowed from the start to the end
	docsControllerPageRegexp = regexp.MustCompile(`\A[a-zA-Z0-9\-/]+\z`)
)

type DocsController struct {
	defaultTmpl string
	errorTmpl   string

	logTag string
	logger boshlog.Logger
}

func NewDocsController(logger boshlog.Logger) DocsController {
	return DocsController{
		defaultTmpl: "index",
		errorTmpl:   "error",

		logTag: "DocsController",
		logger: logger,
	}
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

	r.HTML(200, "docs/"+tmpl, nil)
}
