package release

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/url"

	bprel "github.com/bosh-dep-forks/bosh-provisioner/release"
	"github.com/bosh-io/web/ui/nav"
	semiver "github.com/cppforlife/go-semi-semantic/version"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	bhrelsrepo "github.com/bosh-io/web/release/releasesrepo"
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

	Graph      Graph
	NavPrimary nav.Link

	// memoized notes
	notesInMarkdown *[]byte
}

type Graph interface {
	SVG() template.HTML
}

type releaseAPIRecord struct {
	Name    string `json:"name"`
	Version string `json:"version"`

	URL    string `json:"url"`
	SHA1   string `json:"sha1"`
	SHA256 string `json:"sha256,omitempty"`
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

func (r Release) BuildNavigation(active string) nav.Link {
	root := Navigation()

	relnav := nav.Link{Title: r.Name}
	relnav.Add(nav.Link{
		Title: "All Versions",
		URL:   r.AllVersionsURL(),
	})
	relnav.Add(r.Navigation())
	root.Add(relnav)

	root.Activate(active)

	return root
}

func (r Release) AllURL() string { return "/releases" }

func (r Release) AllVersionsURL() string {
	return fmt.Sprintf("/releases/%s?all=1", r.Source)
}

func (r Release) AvatarURL() string { return r.relVerRec.AvatarURL() }

func (r Release) URL() string {
	return fmt.Sprintf("/releases/%s?version=%s", r.Source, url.QueryEscape(r.Version.AsString()))
}

func (r Release) DownloadURL() string {
	return fmt.Sprintf("/d/%s?v=%s", r.Source, url.QueryEscape(r.Version.AsString()))
}

func (r Release) UserVisibleDownloadURL() string {
	// todo make domain configurable
	return fmt.Sprintf("https://bosh.io/d/%s?v=%s", r.Source, url.QueryEscape(r.Version.AsString()))
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

func (r Release) IsBOSH() bool { return r.Source.IsBOSH() }

func (r Release) IsCPI() bool { return r.Source.IsCPI() }

func (r Release) CPIDocsLink() template.HTML {
	cpi, found := KnownCPIs.FindByShortName(r.Source.ShortName())
	if found {
		return template.HTML(cpi.DocPageLink())
	}

	return template.HTML("")
}

func (r Release) TarballSHA1() (string, error) {
	relTarRec, err := r.relVerRec.Tarball()
	if err != nil {
		return "", err
	}

	return relTarRec.SHA1, nil
}

func (r Release) TarballSHA256() (string, error) {
	relTarRec, err := r.relVerRec.Tarball()
	if err != nil {
		return "", err
	}

	return relTarRec.SHA256, nil
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
	sha1, err := r.TarballSHA1()
	if err != nil {
		return nil, err
	}
	sha256, err := r.TarballSHA256()
	if err != nil {
		return nil, err
	}

	record := releaseAPIRecord{
		Name:    r.Source.Full(),
		Version: r.Version.AsString(),

		URL:    r.UserVisibleDownloadURL(),
		SHA1:   sha1,
		SHA256: sha256,
	}

	return json.Marshal(record)
}

func (r Release) Navigation() nav.Link {
	releaseNav := nav.Link{
		Title: fmt.Sprintf("%s", r.Version),
		URL:   r.URL(),
	}

	releaseNav.Add(nav.Link{
		Title: "Overview",
		URL:   r.URL(),
	})

	{
		jobsNav := nav.Link{Title: "Jobs"}

		for _, job := range r.Jobs {
			jobsNav.Add(nav.Link{
				Title: job.Name,
				URL:   job.URL(),
			})
		}

		releaseNav.Add(jobsNav)
	}

	{
		pkgsNav := nav.Link{Title: "Packages"}

		for _, pkg := range r.Packages {
			pkgsNav.Add(nav.Link{
				Title: pkg.Name,
				URL:   pkg.URL(),
			})
		}

		releaseNav.Add(pkgsNav)
	}

	return releaseNav
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
