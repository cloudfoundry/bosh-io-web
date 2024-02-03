package job_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/bosh-io/web/ui/job"
)

var _ = Describe("NewPropertyItems", func() {
	It("properly organizes items", func() {
		props := []Property{
			Property{
				Name:        "key1",
				Description: "description1",
				Default:     "default1",
			},
			Property{
				Name:        "key2.key1",
				Description: "description2.1",
			},
			Property{
				Name:        "key2.key2",
				Description: "description2.2",
				Default:     "default2.2",
			},
		}

		items := NewPropertyItems(props)

		expectedItems := map[string]*PropertyItem{
			"key1": &PropertyItem{
				FullPath:      "key1",
				Indent:        0,
				Key:           "key1",
				Anchor:        "p=key1",
				MissingValues: false,
				Property:      &props[0],
				Children:      map[string]*PropertyItem{},
			},
			"key2": &PropertyItem{
				FullPath:      "key2",
				Indent:        0,
				Key:           "key2",
				Anchor:        "p=key2",
				MissingValues: false,
				Children: map[string]*PropertyItem{
					"key1": &PropertyItem{
						FullPath:      "key2.key1",
						Indent:        1,
						Key:           "key1",
						Anchor:        "p=key2.key1",
						MissingValues: true,
						Property:      &props[1],
						Children:      map[string]*PropertyItem{},
					},
					"key2": &PropertyItem{
						FullPath:      "key2.key2",
						Indent:        1,
						Key:           "key2",
						Anchor:        "p=key2.key2",
						MissingValues: false,
						Property:      &props[2],
						Children:      map[string]*PropertyItem{},
					},
				},
			},
		}

		Expect(items).To(Equal(expectedItems))
	})
})
