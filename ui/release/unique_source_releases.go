package release

import (
	"sort"

	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

type UniqueSourceReleases []*SameSourceReleases

func NewUniqueSourceReleases(relVerRecs []bhrelsrepo.ReleaseVersionRec) UniqueSourceReleases {
	var uniqueRels UniqueSourceReleases

	bySource := map[string]*SameSourceReleases{}

	for _, relVerRec := range relVerRecs {
		rel := NewIncompleteRelease(relVerRec)

		sameRels, found := bySource[relVerRec.Source]
		if found {
			sameRels.Releases = append(sameRels.Releases, rel)
		} else {
			sameRels := &SameSourceReleases{
				Source:   NewSource(relVerRec.Source),
				Releases: []Release{rel},
			}
			bySource[relVerRec.Source] = sameRels
			uniqueRels = append(uniqueRels, sameRels)
		}
	}

	for _, sameRels := range uniqueRels {
		sort.Sort(sort.Reverse(ReleaseSorting(sameRels.Releases)))
	}

	return uniqueRels
}

func (r UniqueSourceReleases) AllURL() string { return "/releases" }
