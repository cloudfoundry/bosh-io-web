package stemsrepo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestStemsrepo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "stemcell/stemsrepo")
}
