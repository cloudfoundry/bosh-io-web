package metalink

import (
	"sort"
	"strings"
)

func Sort(r *Metalink) {
	sort.Slice(r.Files, func(i, j int) bool {
		return strings.Compare(r.Files[i].Name, r.Files[j].Name) < 0
	})
}
