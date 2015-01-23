package stemcell

import (
	"sort"

	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
)

type UniqueVersionStemcells []*SameVersionStemcells

type UniqueVersionStemcellsSorting UniqueVersionStemcells

func NewUniqueVersionStemcells(ss []bhstemsrepo.Stemcell, limit *int) UniqueVersionStemcells {
	var uniqueStems UniqueVersionStemcells

	byVersion := map[string]*SameVersionStemcells{}

	for _, s := range ss {
		stemcell := NewStemcell(s)
		verStr := stemcell.Version.AsString()

		sameStems, found := byVersion[verStr]
		if found {
			sameStems.Stemcells = append(sameStems.Stemcells, stemcell)
		} else {
			sameStems := &SameVersionStemcells{
				Version:   stemcell.Version,
				Stemcells: []Stemcell{stemcell},
			}
			byVersion[verStr] = sameStems
			uniqueStems = append(uniqueStems, sameStems)
		}
	}

	for _, sameStems := range uniqueStems {
		sort.Sort(StemcellSorting(sameStems.Stemcells))
	}

	sort.Sort(sort.Reverse(UniqueVersionStemcellsSorting(uniqueStems)))

	if limit != nil && len(uniqueStems) > *limit {
		uniqueStems = uniqueStems[0:*limit]
	}

	return uniqueStems
}

func NewLatestVersionStemcells(ss []bhstemsrepo.Stemcell) *SameVersionStemcells {
	uniqStems := NewUniqueVersionStemcells(ss, nil)

	if len(uniqStems) > 0 {
		return uniqStems[0]
	}

	return nil
}

func (s UniqueVersionStemcellsSorting) Len() int           { return len(s) }
func (s UniqueVersionStemcellsSorting) Less(i, j int) bool { return s[i].Version.IsLt(s[j].Version) }
func (s UniqueVersionStemcellsSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
