package release

import (
	"fmt"
	"net/url"
	"sort"

	bprel "github.com/bosh-dep-forks/bosh-provisioner/release"
)

type Package struct {
	Release Release

	Name string

	Fingerprint string
	SHA1        string

	// Package dependencies used at compilation of this package
	Dependencies []Package
}

type PackageSorting []Package

func NewPackages(ps []*bprel.Package, rel Release) []Package {
	pkgs := []Package{}

	for _, p := range ps {
		pkg := Package{
			Release: rel,

			Name: p.Name,

			Fingerprint: p.Fingerprint,
			SHA1:        p.SHA1,

			Dependencies: NewPackages(p.Dependencies, rel),
		}
		pkgs = append(pkgs, pkg)
	}

	sort.Sort(PackageSorting(pkgs))

	return pkgs
}

func (p Package) URL() string {
	return fmt.Sprintf("/packages/%s?source=%s&version=%s", p.Name, p.Release.Source, url.QueryEscape(p.Release.Version.AsString()))
}

func (p Package) HasGithubURL() bool { return p.Release.HasGithubURL() }

func (p Package) GithubURL() string {
	return p.Release.GithubURLForPath("packages/"+p.Name, "")
}

func (p Package) GithubURLOnMaster() string {
	return p.Release.GithubURLForPath("packages/"+p.Name, "master")
}

func (s PackageSorting) Len() int           { return len(s) }
func (s PackageSorting) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s PackageSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
