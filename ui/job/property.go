package job

import (
	"bytes"
	"html/template"
	"sort"
	"strings"

	"github.com/cloudfoundry-incubator/candiedyaml"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	bpreljob "github.com/cppforlife/bosh-provisioner/release/job"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

type Property struct {
	Name        string
	Description string

	Default  interface{}
	Examples []PropertyExample
}

type PropertyExample struct {
	Description string
	Value       interface{}
}

type PropertySorting []Property

func NewProperties(ps []bpreljob.Property) []Property {
	props := []Property{}

	for _, p := range ps {
		props = append(props, NewProperty(p))
	}

	sort.Sort(PropertySorting(props))

	return props
}

func NewProperty(p bpreljob.Property) Property {
	prop := Property{
		Name:        p.Name,
		Description: p.Description,

		Default: p.Default,
	}

	if p.Example != nil {
		prop.Examples = []PropertyExample{
			{Value: p.Example},
		}
	}

	for _, ex := range p.Examples {
		prop.Examples = append(prop.Examples, PropertyExample{
			Description: ex.Description,
			Value:       ex.Value,
		})
	}

	return prop
}

func (p Property) GroupName() string {
	pieces := strings.SplitN(p.Name, ".", 3)

	if len(pieces) > 2 {
		return pieces[0] + "." + pieces[1]
	}

	if len(pieces) == 2 {
		return pieces[0]
	}

	return "."
}

func (p Property) DescriptionInMarkdown() (template.HTML, error) {
	unsafeMarkdown := blackfriday.MarkdownCommon([]byte(p.Description))
	safeMarkdown := bluemonday.UGCPolicy().SanitizeBytes(unsafeMarkdown)

	// todo sanitized markdown
	return template.HTML(safeMarkdown), nil
}

func (p Property) HasDefault() bool {
	return p.Default != nil
}

func (p Property) DefaultAsYAML() (string, error) {
	var b bytes.Buffer

	err := candiedyaml.NewEncoder(&b).Encode(p.Default)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Generating yaml for property '%s' default", p.Name)
	}

	return b.String(), nil
}

func (e PropertyExample) DescriptionInMarkdown() (template.HTML, error) {
	unsafeMarkdown := blackfriday.MarkdownCommon([]byte(e.Description))
	safeMarkdown := bluemonday.UGCPolicy().SanitizeBytes(unsafeMarkdown)

	// todo sanitized markdown
	return template.HTML(safeMarkdown), nil
}

func (e PropertyExample) ValueAsYAML() (string, error) {
	var b bytes.Buffer

	err := candiedyaml.NewEncoder(&b).Encode(e.Value)
	if err != nil {
		return "", bosherr.WrapError(err, "Generating yaml for property example")
	}

	return b.String(), nil
}

func (s PropertySorting) Len() int           { return len(s) }
func (s PropertySorting) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s PropertySorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
