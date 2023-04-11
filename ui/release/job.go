package release

import (
	"fmt"
	"net/url"
	"sort"

	bprel "github.com/bosh-dep-forks/bosh-provisioner/release"
)

type Job struct {
	Release Release

	Name string

	Fingerprint string
	SHA1        string
}

type JobSorting []Job

func NewJobs(js []bprel.Job, rel Release) []Job {
	jobs := []Job{}

	for _, j := range js {
		job := Job{
			Release: rel,

			Name: j.Name,

			Fingerprint: j.Fingerprint,
			SHA1:        j.SHA1,
		}
		jobs = append(jobs, job)
	}

	sort.Sort(JobSorting(jobs))

	return jobs
}

func (j Job) URL() string {
	return fmt.Sprintf("/jobs/%s?source=%s&version=%s", j.Name, j.Release.Source, url.QueryEscape(j.Release.Version.AsString()))
}

func (s JobSorting) Len() int           { return len(s) }
func (s JobSorting) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s JobSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
