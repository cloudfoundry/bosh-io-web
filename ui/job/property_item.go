package job

import (
	"fmt"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type PropertyItem struct {
	Indent int

	FullPath string
	Key      string
	Anchor   string

	// HasDefaults shows if either this item or sub-items have defaults
	MissingValues bool

	// technically could have both...
	Property *Property
	Children map[string]*PropertyItem
}

func NewPropertyItems(props []Property) map[string]*PropertyItem {
	root := PropertyItem{
		Children: map[string]*PropertyItem{},
	}

	for _, prop := range props {
		parts := strings.Split(prop.Name, ".")
		relativeRoot := &root

		for partIdx, part := range parts {
			_, found := relativeRoot.Children[part]
			if !found {
				fullPath := strings.Join(parts[0:partIdx+1], ".")
				relativeRoot.Children[part] = &PropertyItem{
					Indent:   partIdx,
					Key:      part,
					FullPath: fullPath,
					Anchor:   "p=" + fullPath,
					Children: map[string]*PropertyItem{},
				}
			}

			relativeRoot = relativeRoot.Children[part]
		}

		relativeRoot.Property = &prop
		relativeRoot.MissingValues = !prop.HasDefault()
	}

	// TODO un-pointerify or use array like before?

	return root.Children
}

func (i PropertyItem) PoundAnchor() string {
	return "#" + i.Anchor
}

func (i PropertyItem) HasLongKey() bool {
	d, err := i.DefaultAsYAML()
	if err != nil {
		d = ""
	}

	parts := strings.Split(i.Key+": "+d, "\n")

	for _, v := range parts {
		if len(v) > 40 {
			return true
		}
	}

	return false
}

func (i PropertyItem) DefaultAsYAML() (string, error) {
	result, err := i.Property.DefaultAsYAML()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Indenting property '%s' default", i.Property.Name)
	}

	// YAML encoder might add extra new line breaks
	result = strings.Trim(result, "\n")

	parts := strings.Split(result, "\n")

	if len(parts) == 1 {
		return result, nil
	}

	for j, v := range parts {
		// it does mean arrays and hashes will be indented
		parts[j] = fmt.Sprintf("  %s", v)
	}

	return "\n" + strings.Join(parts, "\n"), nil
}
