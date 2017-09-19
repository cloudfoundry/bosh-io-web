package stemcell

import (
	"sort"

	bhstemsrepo "github.com/bosh-io/web/stemcell/stemsrepo"
)

type UniqueNameStemcells []*SameNameStemcells

type SameNameStemcells struct {
	Name      string
	ByVersion UniqueVersionStemcells
}

type UniqueNameStemcellsSorting UniqueNameStemcells

func NewUniqueNameStemcells(ss []bhstemsrepo.Stemcell, filter StemcellFilter) UniqueNameStemcells {
	var result UniqueNameStemcells

	groups := map[string][]bhstemsrepo.Stemcell{}

	for _, s := range ss {
		key := s.Name()
		groups[key] = append(groups[key], s)
	}

	for _, sameSs := range groups {
		sameStems := NewSameNameStemcells(sameSs[0].Name(), sameSs, filter)

		if sameStems.HasAnyStemcells() {
			result = append(result, &sameStems)
		}
	}

	sort.Sort(UniqueNameStemcellsSorting(result))

	return result
}

func NewSameNameStemcells(name string, ss []bhstemsrepo.Stemcell, filter StemcellFilter) SameNameStemcells {
	return SameNameStemcells{
		Name:      name,
		ByVersion: NewUniqueVersionStemcells(ss, filter),
	}
}

func (s UniqueNameStemcells) HasAnyStemcells() bool {
	for _, s_ := range s {
		if s_.HasAnyStemcells() {
			return true
		}
	}

	return false
}

func (s SameNameStemcells) HasAnyStemcells() bool {
	return s.ByVersion.HasAnyStemcells()
}

func (s UniqueNameStemcellsSorting) Len() int           { return len(s) }
func (s UniqueNameStemcellsSorting) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s UniqueNameStemcellsSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
