package job

import (
	"bytes"
	"sort"
	"strings"

	"github.com/cloudfoundry-incubator/candiedyaml"
	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	bpreljob "github.com/cppforlife/bosh-provisioner/release/job"
)

type Property struct {
	Name        string
	Description string

	Default interface{}
	Example interface{}
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
	return Property{
		Name:        p.Name,
		Description: p.Description,

		Default: p.Default,
		Example: p.Example,
	}
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

func (p Property) HasDefault() bool {
	return p.Default != nil
}

func (p Property) DefaultAsYAML() (string, error) {
	var b bytes.Buffer

	err := candiedyaml.NewEncoder(&b).Encode(p.Default)
	if err != nil {
		return "", bosherr.WrapError(err, "Generating yaml for property '%s' default", p.Name)
	}

	return b.String(), nil
}

func (p Property) HasExample() bool {
	return p.Example != nil
}

func (p Property) ExampleAsYAML() (string, error) {
	var b bytes.Buffer

	err := candiedyaml.NewEncoder(&b).Encode(p.Example)
	if err != nil {
		return "", bosherr.WrapError(err, "Generating yaml for property '%s' example", p.Name)
	}

	return b.String(), nil
}

func (s PropertySorting) Len() int           { return len(s) }
func (s PropertySorting) Less(i, j int) bool { return s[i].Name < s[j].Name }
func (s PropertySorting) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
