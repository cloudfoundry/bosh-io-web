package stemcell_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	bhstemsrepo "github.com/bosh-io/web/stemcell/stemsrepo"
	. "github.com/bosh-io/web/ui/stemcell"
)

func s3Stemcell(key string) bhstemsrepo.Stemcell {
	s := bhstemsrepo.NewS3Stemcell(key, "", "", "", 0, "", "")
	Expect(s).ToNot(BeNil())
	return s
}

var _ = Describe("Stemcell", func() {
	DescribeTable("InfrastructureName",
		func(key, expected string) {
			Expect(NewStemcell(s3Stemcell(key)).RegularSource.InfrastructureName()).To(Equal(expected))
		},
		Entry("oci-container", "oci/bosh-stemcell-0.190-oci-container-ubuntu-resolute.tgz", "OCI Container Image"),
		Entry("oci-container rosetta", "oci/bosh-stemcell-0.190-oci-container-ubuntu-resolute-rosetta.tgz", "OCI Container Image (rosetta)"),
		Entry("warden-boshlite", "warden/bosh-stemcell-1.900-warden-boshlite-ubuntu-noble.tgz", "Warden (BOSH Lite)"),
		Entry("gdch-kvm", "gdch/bosh-stemcell-1.900-gdch-kvm-ubuntu-noble.tgz", "Google Distributed Cloud"),
	)
})

var _ = Describe("NewDistroGroups", func() {
	var stemcells []bhstemsrepo.Stemcell

	BeforeEach(func() {
		stemcells = []bhstemsrepo.Stemcell{
			s3Stemcell("warden/bosh-stemcell-0.146-warden-boshlite-ubuntu-resolute.tgz"),
			s3Stemcell("warden/bosh-stemcell-0.93-warden-boshlite-ubuntu-resolute-rosetta.tgz"),
			s3Stemcell("oci/bosh-stemcell-0.190-oci-container-ubuntu-resolute.tgz"),
			s3Stemcell("oci/bosh-stemcell-0.190-oci-container-ubuntu-resolute-rosetta.tgz"),
		}
	})

	names := func(groups DistroGroups) []string {
		var result []string
		for _, g := range groups {
			for _, n := range g.ByName {
				result = append(result, n.Name)
			}
		}
		return result
	}

	It("omits superseded stemcells from the latest versions listing", func() {
		groups := NewDistroGroups(stemcells, StemcellFilter{})

		Expect(names(groups)).To(ConsistOf(
			"bosh-oci-container-ubuntu-resolute",
			"bosh-oci-container-ubuntu-resolute-rosetta",
		))
	})

	It("still lists a superseded stemcell when it is requested by name", func() {
		name := "bosh-warden-boshlite-ubuntu-resolute"
		groups := NewDistroGroups([]bhstemsrepo.Stemcell{stemcells[0]}, StemcellFilter{Name: name})

		Expect(names(groups)).To(ConsistOf(name))
	})
})
