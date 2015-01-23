package job

import (
	"sort"
	"strings"

	bpreljob "github.com/cppforlife/bosh-provisioner/release/job"
)

const (
	templateBinRun       = "bin/run"
	templateBinPrefix    = "bin/"
	templateConfigPrefix = "config/"
)

type Template struct {
	SrcPathEnd string
	DstPathEnd string
}

type TemplateSorting []Template

func NewTemplates(ts []bpreljob.Template) []Template {
	templates := []Template{}

	for _, t := range ts {
		template := Template{
			SrcPathEnd: t.SrcPathEnd,
			DstPathEnd: t.DstPathEnd,
		}
		templates = append(templates, template)
	}

	sort.Sort(TemplateSorting(templates))

	return templates
}

func (t Template) IsBinRun() bool {
	return t.DstPathEnd == templateBinRun
}

func (t Template) IsBin() bool {
	return strings.HasPrefix(t.DstPathEnd, templateBinPrefix)
}

func (t Template) IsConfig() bool {
	return strings.HasPrefix(t.DstPathEnd, templateConfigPrefix)
}

func (t Template) IsOther() bool {
	return !t.IsBin() && !t.IsConfig()
}

func (s TemplateSorting) Len() int           { return len(s) }
func (s TemplateSorting) Less(i, j int) bool { return s[i].DstPathEnd < s[j].DstPathEnd }
func (s TemplateSorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
