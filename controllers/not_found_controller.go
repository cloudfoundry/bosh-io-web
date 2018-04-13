package controllers

import (
	"net/http"
	"os"
)

type NotFoundController struct{}

func (NotFoundController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("templates/docs/404.html")
	if err != nil {
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return
	}

	http.ServeContent(w, r, "templates/docs/404.html", fi.ModTime(), f)
}
