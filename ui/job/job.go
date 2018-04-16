package job

import (
	"fmt"

	bpreljob "github.com/cppforlife/bosh-provisioner/release/job"

	bhrelui "github.com/bosh-io/web/ui/release"
)

type Job struct {
	Release bhrelui.Release

	Name string

	Description string

	Templates []Template
	Packages  []Package

	PropertyItems map[string]*PropertyItem
}

func NewJob(j bpreljob.Job, rel bhrelui.Release) Job {
	props := NewProperties(j.Properties)

	job := Job{
		Release: rel,

		Name: j.Name,

		Description: j.Description,

		Templates: NewTemplates(j.Templates),
		Packages:  NewPackages(j.Packages, rel),

		PropertyItems: NewPropertyItems(props),
	}

	return job
}

func (j Job) URL() string {
	return fmt.Sprintf("/jobs/%s?source=%s&version=%s", j.Name, j.Release.Source, j.Release.Version)
}

func (j Job) HasGithubURL() bool { return j.Release.HasGithubURL() }

func (j Job) GithubURL() string {
	return j.Release.GithubURLForPath("jobs/"+j.Name, "")
}

func (j Job) GithubURLOnMaster() string {
	return j.Release.GithubURLForPath("jobs/"+j.Name, "master")
}

func (j Job) IsErrand() bool {
	ts := j.selectTemplates(func(t Template) bool { return t.IsBinRun() })
	return len(ts) > 0
}

func (j Job) BinTemplates() []Template {
	return j.selectTemplates(func(t Template) bool { return t.IsBin() })
}

func (j Job) ConfigTemplates() []Template {
	return j.selectTemplates(func(t Template) bool { return t.IsConfig() })
}

func (j Job) OtherTemplates() []Template {
	return j.selectTemplates(func(t Template) bool { return t.IsOther() })
}

func (j Job) selectTemplates(pred func(t Template) bool) []Template {
	templates := []Template{}

	for _, t := range j.Templates {
		if pred(t) {
			templates = append(templates, t)
		}
	}

	return templates
}
