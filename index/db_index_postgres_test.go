package index_test

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cppforlife/bosh-hub/index"
)

var _ = Describe("DBIndex with postgres", func() {
	var (
		index DBIndex
	)

	BeforeEach(func() {
		url := "postgres://postgres@localhost/pivotal?sslmode=disable"

		logger := boshlog.NewLogger(boshlog.LevelNone)

		adapterPool, err := NewPostgresAdapterPool(url, logger)
		Expect(err).ToNot(HaveOccurred())

		adapter, err := adapterPool.NewAdapter("pivotal")
		Expect(err).ToNot(HaveOccurred())

		_, err = adapter.Clear()
		Expect(err).ToNot(HaveOccurred())

		index = NewDBIndex(adapter, logger)
	})

	ItBehavesLikeAnIndex(&index)
})
