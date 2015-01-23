package release

import (
	"fmt"
	"sort"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

type SameSourceReleases struct {
	Source   Source
	Releases []Release
}

func NewSameSourceReleases(source string, relVerRecs []bhrelsrepo.ReleaseVersionRec) SameSourceReleases {
	rels := SameSourceReleases{
		Source: NewSource(source),
	}

	for _, relVerRec := range relVerRecs {
		rel := NewIncompleteRelease(source, relVerRec.Version())
		rels.Releases = append(rels.Releases, rel)
	}

	sort.Sort(sort.Reverse(ReleaseSorting(rels.Releases)))

	if len(rels.Releases) > 0 {
		rels.Releases[0].IsLatest = true
	}

	return rels
}

func (r SameSourceReleases) FirstXReleases(x int) []Release {
	if len(r.Releases) < x {
		return r.Releases
	}

	return r.Releases[0:x]
}

func (r SameSourceReleases) HasMoreThanXReleases(x int) bool {
	return len(r.Releases) > x
}

func (r SameSourceReleases) AllURL() string { return "/releases" }

func (r SameSourceReleases) URL() string {
	return fmt.Sprintf("/releases/%s", r.Source)
}
