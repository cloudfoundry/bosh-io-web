package stemcell

import (
	"sort"
)

type DistroGroup struct {
	Distro Distro

	Stemcells []Stemcell
}

type DistroGroups []DistroGroup

type DistroGroupSorting []DistroGroup

func NewDistroGroups(stemcells []Stemcell) DistroGroups {
	var groups []DistroGroup

	for _, d := range supportedDistros {
		groups = append(groups, DistroGroup{Distro: d})
	}

	// Catch any other stemcell
	groups = append(groups, DistroGroup{Distro: unknownDistro})

	for _, stemcell := range stemcells {
		for i, g := range groups {
			if g.Distro.Matches(stemcell) {
				groups[i].Stemcells = append(groups[i].Stemcells, stemcell)
				break
			}
		}
	}

	sort.Sort(sort.Reverse(DistroGroupSorting(groups)))

	return groups
}

func (gs DistroGroups) AllStemcellsLen() int {
	result := 0

	for _, g := range gs {
		result += len(g.Stemcells)
	}

	return result
}

func (s DistroGroupSorting) Len() int           { return len(s) }
func (s DistroGroupSorting) Less(i, j int) bool { return len(s[i].Stemcells) < len(s[j].Stemcells) }
func (s DistroGroupSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
