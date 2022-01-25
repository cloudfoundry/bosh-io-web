package stemsrepo_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/bosh-io/web/stemcell/stemsrepo"
)

var _ = Describe("NewS3Stemcell", func() {
	type ExtractedPieces struct {
		Version string
		Name    string

		InfName    string
		HvName     string
		DiskFormat string

		OSName    string
		OSVersion string

		AgentType string
	}

	var examples = map[string]ExtractedPieces{
		"bosh-stemcell/aws/bosh-stemcell-891-aws-xen-ubuntu.tgz": ExtractedPieces{
			Name:    "bosh-aws-xen-ubuntu",
			Version: "891",

			InfName: "aws",
			HvName:  "xen",

			OSName:    "ubuntu",
			OSVersion: "lucid",

			AgentType: "ruby",
		},

		"bosh-stemcell/aws/bosh-stemcell-2311-aws-xen-centos-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-aws-xen-centos-go_agent",
			Version: "2311",

			InfName: "aws",
			HvName:  "xen",

			OSName:    "centos",
			OSVersion: "",

			AgentType: "go",
		},

		"bosh-stemcell/aws/bosh-stemcell-2446-aws-xen-ubuntu-lucid-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-aws-xen-ubuntu-lucid-go_agent",
			Version: "2446",

			InfName: "aws",
			HvName:  "xen",

			OSName:    "ubuntu",
			OSVersion: "lucid",

			AgentType: "go",
		},

		"micro-bosh-stemcell/aws/light-micro-bosh-stemcell-891-aws-xen-ubuntu.tgz": ExtractedPieces{
			Name:    "bosh-aws-xen-ubuntu",
			Version: "891",

			InfName: "aws",
			HvName:  "xen",

			OSName:    "ubuntu",
			OSVersion: "lucid",

			AgentType: "ruby",
		},

		"micro-bosh-stemcell/warden/bosh-stemcell-56-warden-boshlite-ubuntu-lucid-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-warden-boshlite-ubuntu-lucid-go_agent",
			Version: "56",

			InfName: "warden",
			HvName:  "boshlite",

			OSName:    "ubuntu",
			OSVersion: "lucid",

			AgentType: "go",
		},

		"bosh-stemcell/aws/light-bosh-stemcell-2579-aws-xen-centos.tgz": ExtractedPieces{
			Name:    "bosh-aws-xen-centos",
			Version: "2579",

			InfName: "aws",
			HvName:  "xen",

			OSName:    "centos",
			OSVersion: "",

			AgentType: "ruby",
		},

		"bosh-stemcell/aws/light-bosh-stemcell-2579-aws-xen-centos-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-aws-xen-centos-go_agent",
			Version: "2579",

			InfName: "aws",
			HvName:  "xen",

			OSName:    "centos",
			OSVersion: "",

			AgentType: "go",
		},

		"bosh-stemcell/aws/light-bosh-stemcell-2579.1-aws-xen-centos-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-aws-xen-centos-go_agent",
			Version: "2579.1",

			InfName: "aws",
			HvName:  "xen",

			OSName:    "centos",
			OSVersion: "",

			AgentType: "go",
		},

		"bosh-stemcell/aws/light-bosh-stemcell-2579.1-aws-xen-hvm-centos-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-aws-xen-hvm-centos-go_agent",
			Version: "2579.1",

			InfName: "aws",
			HvName:  "xen-hvm",

			OSName:    "centos",
			OSVersion: "",

			AgentType: "go",
		},

		"bosh-stemcell/aws/light-bosh-stemcell-2579.1-aws-xen-hvm-ubuntu-trusty-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-aws-xen-hvm-ubuntu-trusty-go_agent",
			Version: "2579.1",

			InfName: "aws",
			HvName:  "xen-hvm",

			OSName:    "ubuntu",
			OSVersion: "trusty",

			AgentType: "go",
		},

		// Notice no top-level folder prefix
		"aws/bosh-stemcell-3306-aws-xen-ubuntu-trusty-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-aws-xen-ubuntu-trusty-go_agent",
			Version: "3306",

			InfName: "aws",
			HvName:  "xen",

			OSName:    "ubuntu",
			OSVersion: "trusty",

			AgentType: "go",
		},

		// Notice no folder prefixes
		"bosh-stemcell-2776-warden-boshlite-centos-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-warden-boshlite-centos-go_agent",
			Version: "2776",

			InfName: "warden",
			HvName:  "boshlite",

			OSName:    "centos",
			OSVersion: "",

			AgentType: "go",
		},

		// Disk format
		"bosh-stemcell/openstack/bosh-stemcell-56-openstack-kvm-ubuntu-trusty-go_agent-raw.tgz": ExtractedPieces{
			Name:    "bosh-openstack-kvm-ubuntu-trusty-go_agent-raw",
			Version: "56",

			InfName:    "openstack",
			HvName:     "kvm",
			DiskFormat: "raw",

			OSName:    "ubuntu",
			OSVersion: "trusty",

			AgentType: "go",
		},

		// Numeric OS version
		"bosh-stemcell/vsphere/bosh-stemcell-2922-vsphere-esxi-centos-7-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-vsphere-esxi-centos-7-go_agent",
			Version: "2922",

			InfName:    "vsphere",
			HvName:     "esxi",
			DiskFormat: "",

			OSName:    "centos",
			OSVersion: "7",

			AgentType: "go",
		},

		// Stemcell in China
		"light-china-bosh-stemcell-3130-aws-xen-hvm-ubuntu-trusty-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-aws-xen-hvm-ubuntu-trusty-go_agent",
			Version: "3130",

			InfName: "aws",
			HvName:  "xen-hvm",

			OSName:    "ubuntu",
			OSVersion: "trusty",

			AgentType: "go",
		},

		"light-bosh-stemcell-1709.3-google-kvm-windows2016-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-google-kvm-windows2016-go_agent",
			Version: "1709.3",

			InfName: "google",
			HvName:  "kvm",

			OSName:    "windows",
			OSVersion: "2016",

			AgentType: "go",
		},

		"light-bosh-stemcell-1089.0-aws-xen-hvm-windows2012R2-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-aws-xen-hvm-windows2012R2-go_agent",
			Version: "1089.0",

			InfName: "aws",
			HvName:  "xen-hvm",

			OSName:    "windows",
			OSVersion: "2012R2",

			AgentType: "go",
		},

		"light-bosh-stemcell-1089.0-azure-hyperv-windows2012R2-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-azure-hyperv-windows2012R2-go_agent",
			Version: "1089.0",

			InfName: "azure",
			HvName:  "hyperv",

			OSName:    "windows",
			OSVersion: "2012R2",

			AgentType: "go",
		},

		"light-bosh-stemcell-1089.0-google-kvm-windows2012R2-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-google-kvm-windows2012R2-go_agent",
			Version: "1089.0",

			InfName: "google",
			HvName:  "kvm",

			OSName:    "windows",
			OSVersion: "2012R2",

			AgentType: "go",
		},

		// Softlayer stemcell
		"light-bosh-stemcell-3232.4-softlayer-esxi-ubuntu-trusty-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-softlayer-esxi-ubuntu-trusty-go_agent",
			Version: "3232.4",

			InfName: "softlayer",
			HvName:  "esxi",

			OSName:    "ubuntu",
			OSVersion: "trusty",

			AgentType: "go",
		},

		// Ubuntu xenial
		"azure/bosh-stemcell-40-azure-hyperv-ubuntu-xenial-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-azure-hyperv-ubuntu-xenial-go_agent",
			Version: "40",

			InfName: "azure",
			HvName:  "hyperv",

			OSName:    "ubuntu",
			OSVersion: "xenial",

			AgentType: "go",
		},

		// Ubuntu bionic
		"azure/bosh-stemcell-40-azure-hyperv-ubuntu-bionic-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-azure-hyperv-ubuntu-bionic-go_agent",
			Version: "40",

			InfName: "azure",
			HvName:  "hyperv",

			OSName:    "ubuntu",
			OSVersion: "bionic",

			AgentType: "go",
		},

		// Ubuntu jammy
		"azure/bosh-stemcell-40-azure-hyperv-ubuntu-jammy-go_agent.tgz": ExtractedPieces{
			Name:    "bosh-azure-hyperv-ubuntu-jammy-go_agent",
			Version: "40",

			InfName: "azure",
			HvName:  "hyperv",

			OSName:    "ubuntu",
			OSVersion: "jammy",

			AgentType: "go",
		},
	}

	for p, e := range examples {
		path := p
		example := e

		It(fmt.Sprintf("correctly interprets '%s'", path), func() {
			s3Stemcell := NewS3Stemcell(path, "", "", "", 0, "", "")
			Expect(s3Stemcell).ToNot(BeNil())

			Expect(s3Stemcell.Name()).To(Equal(example.Name))
			Expect(s3Stemcell.Version().AsString()).To(Equal(example.Version))

			Expect(s3Stemcell.InfName()).To(Equal(example.InfName))
			Expect(s3Stemcell.HvName()).To(Equal(example.HvName))
			Expect(s3Stemcell.DiskFormat()).To(Equal(example.DiskFormat))

			Expect(s3Stemcell.OSName()).To(Equal(example.OSName))
			Expect(s3Stemcell.OSVersion()).To(Equal(example.OSVersion))

			Expect(s3Stemcell.AgentType()).To(Equal(example.AgentType))
		})
	}
})
