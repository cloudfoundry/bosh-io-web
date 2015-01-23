package job_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cppforlife/bosh-hub/job"
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

		expectedItems := []PropertyItem{
			PropertyItem{
				Indent:        0,
				Key:           "key1",
				MissingValues: false,
				Property:      &props[0],
			},
			PropertyItem{
				Indent:        0,
				Key:           "key2",
				MissingValues: true,
			},
			PropertyItem{
				Indent:        1,
				Key:           "key1",
				MissingValues: true,
				Property:      &props[1],
			},
			PropertyItem{
				Indent:        1,
				Key:           "key2",
				MissingValues: false,
				Property:      &props[2],
			},
		}

		Expect(items).To(Equal(expectedItems))
	})
})
