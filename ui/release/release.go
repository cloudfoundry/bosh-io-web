package release

import (
	"encoding/json"
	"fmt"
	"html/template"

	bprel "github.com/cppforlife/bosh-provisioner/release"
	semiver "github.com/cppforlife/go-semi-semantic/version"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

type Release struct {
	relVerRec bhrelsrepo.ReleaseVersionRec

	Source Source

	Name    string
	Version semiver.Version

	IsLatest bool

	CommitHash string

	Jobs []Job

	Packages []Package

	Graph Graph

	// memoized notes
	notesInMarkdown *[]byte
}

type Graph interface {
	SVG() template.HTML
}

type releaseAPIRecord struct {
	Name    string `json:"name"`
	Version string `json:"version"`

	URL string `json:"url"`
}

type ReleaseSorting []Release

func NewRelease(relVerRec bhrelsrepo.ReleaseVersionRec, r bprel.Release) Release {
	rel := Release{
		relVerRec: relVerRec,

		Source: NewSource(relVerRec.AsSource()),

		Name:    r.Name,
		Version: relVerRec.Version(),

		CommitHash: r.CommitHash,

		IsLatest: false,
	}

	rel.Jobs = NewJobs(r.Jobs, rel)
	rel.Packages = NewPackages(r.Packages, rel)

	return rel
}

func NewIncompleteRelease(relVerRec bhrelsrepo.ReleaseVersionRec, name string) Release {
	return Release{
		relVerRec: relVerRec,

		Source: NewSource(relVerRec.AsSource()),

		Name:    name,
		Version: relVerRec.Version(),
	}
}

func (r Release) AllURL() string { return "/releases" }

func (r Release) AllVersionsURL() string {
	return fmt.Sprintf("/releases/%s", r.Source)
}

func (r Release) AvatarURL() string { return r.relVerRec.AvatarURL() }

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

func (r Release) IsCPI() bool { return r.Source.IsCPI() }

func (r Release) CPIDocsLink() template.HTML {
	links := map[string][]string{
		"bosh-aws-cpi-release":       []string{"AWS", "/docs/init-aws.html"},
		"bosh-openstack-cpi-release": []string{"OpenStack", "/docs/init-openstack.html"},
		"bosh-vsphere-cpi-release":   []string{"vSphere", "/docs/init-vsphere.html"},
		"bosh-vcloud-cpi-release":    []string{"vCloud", "/docs/init-vcloud.html"},
	}

	if link, found := links[r.Source.ShortName()]; found {
		return template.HTML(fmt.Sprintf("<a href='%s'>Initializing a BOSH environment on %s</a>", link[1], link[0]))
	}

	return template.HTML("")
}

func (r Release) TarballSHA1() (string, error) {
	relTarRec, _, err := r.relVerRec.Tarball()
	if err != nil {
		return "", err
	}

	return relTarRec.SHA1, nil
}

func (r *Release) NotesInMarkdown() (template.HTML, error) {
	if r.notesInMarkdown == nil {
		// Do not care about found -> no UI indicator
		noteRec, _, err := r.relVerRec.Notes()
		if err != nil {
			return template.HTML(""), err
		}

		unsafeMarkdown := blackfriday.MarkdownCommon([]byte(noteRec.Content))
		safeMarkdown := bluemonday.UGCPolicy().SanitizeBytes(unsafeMarkdown)

		r.notesInMarkdown = &safeMarkdown
	}

	// todo sanitized markdown
	return template.HTML(*r.notesInMarkdown), nil
}

func (r Release) MarshalJSON() ([]byte, error) {
	record := releaseAPIRecord{
		Name:    r.Source.Full(),
		Version: r.Version.AsString(),

		// todo MD5 is expensive to get
		URL: r.UserVisibleDownloadURL(),
	}

	return json.Marshal(record)
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
