package release

import (
	"fmt"
	"html/template"
	"strings"

	bprel "github.com/cppforlife/bosh-provisioner/release"
	semiver "github.com/cppforlife/go-semi-semantic/version"
)

type Release struct {
	Source Source

	Name    string
	Version semiver.Version

	IsLatest bool

	CommitHash string

	Jobs []Job

	Packages []Package

	Graph Graph
}

type Graph interface {
	SVG() template.HTML
}

type ReleaseSorting []Release

func NewRelease(source string, r bprel.Release) Release {
	rel := Release{
		Source: NewSource(source),

		Name:    r.Name,
		Version: parseVersion(r.Version),

		CommitHash: r.CommitHash,

		IsLatest: false,
	}

	rel.Jobs = NewJobs(r.Jobs, rel)
	rel.Packages = NewPackages(r.Packages, rel)

	return rel
}

func NewIncompleteRelease(source string, version semiver.Version) Release {
	return Release{
		Source:  NewSource(source),
		Version: version,
	}
}

func (r Release) AllURL() string { return "/releases" }

func (r Release) AllVersionsURL() string {
	return fmt.Sprintf("/releases/%s", r.Source)
}

func (r Release) URL() string {
	return fmt.Sprintf("/releases/%s?version=%s", r.Source, r.Version)
}

func (r Release) DownloadURL() string {
	return fmt.Sprintf("/d/%s?v=%s", r.Source, r.Version)
}

func (r Release) UserVisibleDownloadURL() string {
	// todo make domain configurable
	return fmt.Sprintf("https://bosh.io/d/%s?v=%s", r.Source, r.Version)
}

func (r Release) UserVisibleLatestDownloadURL() string {
	// todo make domain configurable
	return fmt.Sprintf("https://bosh.io/d/%s", r.Source)
}

func (r Release) GraphURL() string { return r.URL() + "&graph=1" }

func (r Release) HasGithubURL() bool { return r.Source.FromGithub() }

func (r Release) GithubURL() string {
	return r.GithubURLForPath("", "")
}

func (r Release) GithubURLOnMaster() string {
	return r.GithubURLForPath("", "master")
}

func (r Release) GithubURLForPath(path, ref string) string {
	if len(ref) > 0 {
		// nothing
	} else if len(r.CommitHash) > 0 {
		ref = r.CommitHash
	} else {
		// Some releases might not have CommitHash
		ref = "<missing>"
	}

	// e.g. https://github.com/cloudfoundry/cf-release/tree/1c96107/jobs/hm9000
	return fmt.Sprintf("%s/tree/%s/%s", r.Source.GithubURL(), ref, path)
}

func (r Release) IsCPI() bool {
	return strings.HasSuffix(r.Name, "-cpi")
}

func (s ReleaseSorting) Len() int           { return len(s) }
func (s ReleaseSorting) Less(i, j int) bool { return s[i].Version.IsLt(s[j].Version) }
func (s ReleaseSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func parseVersion(versionRaw string) semiver.Version {
	ver, err := semiver.NewVersionFromString(versionRaw)
	if err != nil {
		panic(fmt.Sprintf("Version '%s' is not valid: %s", versionRaw, err))
	}

	return ver
}
