package release

import (
	"sort"

	bhrelsrepo "github.com/bosh-io/web/release/releasesrepo"
)

type UniqueSources []Source

func NewUniqueSources(srcs []bhrelsrepo.Source) UniqueSources {
	var sources UniqueSources

	for _, src := range srcs {
		sources = append(sources, NewSource(src))
	}

	sort.Sort(SourceSorting(sources))

	return sources
}
