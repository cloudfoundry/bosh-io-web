package controllers

import "net/http"

type RedirectsConfig map[string]string

type RedirectsController struct {
	redirects RedirectsConfig
}

func NewRedirectsController(
	redirects RedirectsConfig,
) RedirectsController {
	return RedirectsController{
		redirects: redirects,
	}
}

func (c RedirectsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	target, found := c.redirects[r.URL.Path]
	if !found {
		return
	}

	http.Redirect(w, r, target, http.StatusMovedPermanently)
}
