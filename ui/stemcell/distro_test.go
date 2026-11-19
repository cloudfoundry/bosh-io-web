package stemcell_test

import (
	. "github.com/bosh-io/web/ui/stemcell"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("IsListedDistro", func() {
	DescribeTable("reports whether a stemcell belongs to a listed distro",
		func(key string, expected bool) {
			Expect(IsListedDistro(s3Stemcell(key))).To(Equal(expected))
		},
		Entry("ubuntu-jammy on aws", "light-bosh-stemcell-1.900-aws-xen-hvm-ubuntu-jammy-go_agent.tgz", true),
		Entry("ubuntu-noble on google", "light-bosh-stemcell-1.900-google-kvm-ubuntu-noble-go_agent.tgz", true),
		Entry("ubuntu-resolute on aws", "light-bosh-stemcell-1.900-aws-xen-hvm-ubuntu-resolute-go_agent.tgz", true),
		Entry("windows2019 on aws", "light-bosh-stemcell-2019.105-aws-xen-hvm-windows2019-go_agent.tgz", true),
		Entry("ubuntu-jammy on softlayer", "bosh-stemcell-1.900-softlayer-xen-ubuntu-jammy-go_agent.tgz", true),
		Entry("ubuntu-jammy on alicloud", "bosh-stemcell-1.900-alicloud-kvm-ubuntu-jammy-go_agent.tgz", true),
		Entry("ubuntu-resolute-rosetta on warden", "warden/bosh-stemcell-0.181-warden-boshlite-ubuntu-resolute-rosetta-go_agent.tgz", true),
		Entry("ubuntu-xenial", "light-bosh-stemcell-621.125-aws-xen-hvm-ubuntu-xenial-go_agent.tgz", false),
		Entry("ubuntu-trusty", "light-bosh-stemcell-3586.100-aws-xen-hvm-ubuntu-trusty-go_agent.tgz", false),
		Entry("ubuntu-bionic", "light-bosh-stemcell-1.204-aws-xen-hvm-ubuntu-bionic-go_agent.tgz", false),
		Entry("centos-7", "light-bosh-stemcell-3763.61-aws-xen-hvm-centos-7-go_agent.tgz", false),
		Entry("windows2016", "light-bosh-stemcell-1709.21-aws-xen-hvm-windows2016-go_agent.tgz", false),
		Entry("windows2012R2", "light-bosh-stemcell-1200.38-aws-xen-hvm-windows2012R2-go_agent.tgz", false),
		Entry("windows1803", "light-bosh-stemcell-1803.17-aws-xen-hvm-windows1803-go_agent.tgz", false),
	)
})
