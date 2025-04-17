package stemcell

import (
	"sort"

	semiver "github.com/cppforlife/go-semi-semantic/version"

	bhstemsrepo "github.com/bosh-io/web/stemcell/stemsrepo"
)

type UniqueVersionStemcells []*SameVersionStemcells

type SameVersionStemcells struct {
	Version   semiver.Version
	Stemcells []Stemcell

	ShowingAllVersions bool
}

type UniqueVersionStemcellsSorting UniqueVersionStemcells

func NewUniqueVersionStemcells(ss []bhstemsrepo.Stemcell, filter StemcellFilter) UniqueVersionStemcells {
	var result UniqueVersionStemcells

	groups := map[string][]bhstemsrepo.Stemcell{}

	for _, s := range ss {
		key := s.Version().AsString()
		groups[key] = append(groups[key], s)
	}

	for _, sameSs := range groups {
		sameStems := NewSameVersionStemcells(sameSs[0].Version(), sameSs)
		sameStems.ShowingAllVersions = filter.ShowingAllVersions()
		result = append(result, &sameStems)
	}

	sort.Sort(sort.Reverse(UniqueVersionStemcellsSorting(result)))

	if filter.HasLimit() && len(result) > filter.Limit() {
		result = result[0:filter.Limit()]
	}

	return result
}

func NewSameVersionStemcells(version semiver.Version, ss []bhstemsrepo.Stemcell) SameVersionStemcells {
	var stemcells []Stemcell

	groups := map[string]*Stemcell{}

	// Collapse light and non-light stemcells into a single Stemcell UI
	for _, s := range ss {
		if stemcell, ok := groups[s.Name()]; ok {
			stemcell.AddAsSource(s)
		} else {
			stemcell := NewStemcell(s)
			groups[s.Name()] = &stemcell
		}
	}

	for _, stemcell := range groups {
		stemcells = append(stemcells, *stemcell)
	}

	sort.Sort(StemcellManifestNameSorting(stemcells))

	return SameVersionStemcells{Version: version, Stemcells: stemcells}
}

func (s UniqueVersionStemcells) HasAnyStemcells() bool {
	for _, s_ := range s {
		if len(s_.Stemcells) > 0 {
			return true
		}
	}

	return false
}

func (s UniqueVersionStemcells) ForAPI() []Stemcell {
	// API should return empty array, not null!
	stemcells := []Stemcell{}

	for _, uniqueStems := range s {
		for _, stem := range uniqueStems.Stemcells { //nolint:staticcheck
			stemcells = append(stemcells, stem)
		}
	}

	return stemcells
}

func (s SameVersionStemcells) AllURL() string { return "/stemcells" }

func (s UniqueVersionStemcellsSorting) Len() int           { return len(s) }
func (s UniqueVersionStemcellsSorting) Less(i, j int) bool { return s[i].Version.IsLt(s[j].Version) }
func (s UniqueVersionStemcellsSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
