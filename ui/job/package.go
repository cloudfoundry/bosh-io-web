package job

import (
	"fmt"
	"sort"

	bpreljob "github.com/bosh-dep-forks/bosh-provisioner/release/job"

	bhrelui "github.com/bosh-io/web/ui/release"
)

type Package struct {
	Release bhrelui.Release

	Name string
}

type PackageSorting []Package

func NewPackages(ps []bpreljob.Package, rel bhrelui.Release) []Package {
	pkgs := []Package{}

	for _, p := range ps {
		pkg := Package{
			Release: rel,

			Name: p.Name,
		}
		pkgs = append(pkgs, pkg)
	}

	sort.Sort(PackageSorting(pkgs))

	return pkgs
}

func (p Package) URL() string {
	return fmt.Sprintf("/packages/%s?source=%s&version=%s", p.Name, p.Release.Source, p.Release.Version)
}

func (s PackageSorting) Len() int           { return len(s) }
func (s PackageSorting) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s PackageSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
