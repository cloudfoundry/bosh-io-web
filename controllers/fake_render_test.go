package controllers_test

import (
	"html/template"
	"net/http"

	martrend "github.com/martini-contrib/render"
)

type FakeRender struct {
	JSONStatus   int
	JSONResponse interface{}

	responded bool
}

func (r *FakeRender) JSON(status int, v interface{}) {
	if r.responded {
		panic("Responding more than once")
	}
	r.JSONStatus = status
	r.JSONResponse = v
	r.responded = true
}

func (r *FakeRender) HTML(status int, name string, v interface{}, htmlOpt ...martrend.HTMLOptions) {}

func (r *FakeRender) XML(status int, v interface{})           {}
func (r *FakeRender) Data(status int, v []byte)               {}
func (r *FakeRender) Text(status int, v string)               {}
func (r *FakeRender) Error(status int)                        {}
func (r *FakeRender) Status(status int)                       {}
func (r *FakeRender) Redirect(location string, status ...int) {}

func (r *FakeRender) Template() *template.Template { return nil }
func (r *FakeRender) Header() http.Header          { return http.Header{} }
