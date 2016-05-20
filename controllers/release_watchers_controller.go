package controllers

import (
	"net/http"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	mart "github.com/go-martini/martini"
	martrend "github.com/martini-contrib/render"

	bhwatchersrepo "github.com/cppforlife/bosh-hub/release/watchersrepo"
)

type ReleaseWatchersController struct {
	watchersRepo bhwatchersrepo.WatchersRepository
	privateToken string

	indexTmpl string
	errorTmpl string

	logger boshlog.Logger
}

func NewReleaseWatchersController(
	watchersRepo bhwatchersrepo.WatchersRepository,
	privateToken string,
	logger boshlog.Logger,
) ReleaseWatchersController {
	return ReleaseWatchersController{
		watchersRepo: watchersRepo,
		privateToken: privateToken,

		indexTmpl: "release_watchers/index",
		errorTmpl: "error",

		logger: logger,
	}
}

type releaseWatchersControllerIndexPage struct {
	Watchers     []bhwatchersrepo.WatcherRec
	PrivateToken string
}

func (c ReleaseWatchersController) Index(r martrend.Render) {
	watcherRecs, err := c.watchersRepo.ListAll()
	if err != nil {
		r.HTML(500, c.errorTmpl, err)
		return
	}

	page := releaseWatchersControllerIndexPage{
		Watchers:     watcherRecs,
		PrivateToken: c.privateToken,
	}

	r.HTML(200, c.indexTmpl, page)
}

func (c ReleaseWatchersController) WatchOrUnwatch(req *http.Request, r martrend.Render, params mart.Params) {
	action := req.FormValue("action")

	switch {
	case action == "watch":
		err := c.watchersRepo.Add(req.FormValue("relSource"), req.FormValue("minVersion"))
		if err != nil {
			r.HTML(500, c.errorTmpl, err)
			return
		}

	case action == "unwatch":
		err := c.watchersRepo.Remove(req.FormValue("relSource"))
		if err != nil {
			r.HTML(500, c.errorTmpl, err)
			return
		}

	default:
		r.HTML(500, c.errorTmpl, bosherr.Errorf("Unknown action '%s'", action))
		return
	}

	r.Redirect(c.privateToken + "/release_watchers")
}
