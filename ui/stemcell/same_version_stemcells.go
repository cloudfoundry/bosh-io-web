package stemcell

import (
	"sort"

	semiver "github.com/cppforlife/go-semi-semantic/version"

	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
)

type SameVersionStemcells struct {
	Version   semiver.Version
	Stemcells []Stemcell
}

func NewSameVersionStemcells(version semiver.Version, ss []bhstemsrepo.Stemcell) SameVersionStemcells {
	result := SameVersionStemcells{
		Version: version,
	}

	for _, s := range ss {
		stemcell := NewStemcell(s)
		result.Stemcells = append(result.Stemcells, stemcell)
	}

	sort.Sort(sort.Reverse(StemcellSorting(result.Stemcells)))

	return result
}

func (s SameVersionStemcells) DistroGroups() DistroGroups {
	return NewDistroGroups(s.Stemcells)
}

func (s SameVersionStemcells) AllURL() string { return "/stemcells" }
