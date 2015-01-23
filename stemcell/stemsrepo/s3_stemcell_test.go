package stemsrepo_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
)

var _ = Describe("NewS3Stemcell", func() {
	type ExtractedPieces struct {
		Version string
		Name    string

		InfName string
		HvName  string

		OSName    string
		OSVersion string

		AgentType string
	}

	var examples = map[string]ExtractedPieces{
		"bosh-stemcell/aws/bosh-stemcell-891-aws-xen-ubuntu.tgz": ExtractedPieces{
			Name:    "bosh",
			Version: "891",

			InfName: "aws",
			HvName:  "xen",

			OSName:    "ubuntu",
			OSVersion: "lucid",

			AgentType: "ruby_agent",
		},

		"bosh-stemcell/aws/bosh-stemcell-2311-aws-xen-centos-go_agent.tgz": ExtractedPieces{
			Name:    "bosh",
			Version: "2311",

			InfName: "aws",
			HvName:  "xen",

			OSName:    "centos",
			OSVersion: "",

			AgentType: "go_agent",
		},

		"bosh-stemcell/aws/bosh-stemcell-2446-aws-xen-ubuntu-lucid-go_agent.tgz": ExtractedPieces{
			Name:    "bosh",
			Version: "2446",

			InfName: "aws",
			HvName:  "xen",

			OSName:    "ubuntu",
			OSVersion: "lucid",

			AgentType: "go_agent",
		},

		"micro-bosh-stemcell/aws/light-micro-bosh-stemcell-891-aws-xen-ubuntu.tgz": ExtractedPieces{
			Name:    "light-micro-bosh",
			Version: "891",

			InfName: "aws",
			HvName:  "xen",

			OSName:    "ubuntu",
			OSVersion: "lucid",

			AgentType: "ruby_agent",
		},

		"micro-bosh-stemcell/warden/bosh-stemcell-56-warden-boshlite-ubuntu-lucid-go_agent.tgz": ExtractedPieces{
			Name:    "bosh",
			Version: "56",

			InfName: "warden",
			HvName:  "boshlite",

			OSName:    "ubuntu",
			OSVersion: "lucid",

			AgentType: "go_agent",
		},

		"bosh-stemcell/aws/light-bosh-stemcell-2579-aws-xen-centos.tgz": ExtractedPieces{
			Name:    "light-bosh",
			Version: "2579",

			InfName: "aws",
			HvName:  "xen",

			OSName:    "centos",
			OSVersion: "",

			AgentType: "ruby_agent",
		},

		"bosh-stemcell/aws/light-bosh-stemcell-2579-aws-xen-centos-go_agent.tgz": ExtractedPieces{
			Name:    "light-bosh",
			Version: "2579",

			InfName: "aws",
			HvName:  "xen",

			OSName:    "centos",
			OSVersion: "",

			AgentType: "go_agent",
		},

		"bosh-stemcell/aws/light-bosh-stemcell-2579.1-aws-xen-centos-go_agent.tgz": ExtractedPieces{
			Name:    "light-bosh",
			Version: "2579.1",

			InfName: "aws",
			HvName:  "xen",

			OSName:    "centos",
			OSVersion: "",

			AgentType: "go_agent",
		},

		"bosh-stemcell/aws/light-bosh-stemcell-2579.1-aws-xen-hvm-centos-go_agent.tgz": ExtractedPieces{
			Name:    "light-bosh",
			Version: "2579.1",

			InfName: "aws",
			HvName:  "xen-hvm",

			OSName:    "centos",
			OSVersion: "",

			AgentType: "go_agent",
		},

		"bosh-stemcell/aws/light-bosh-stemcell-2579.1-aws-xen-hvm-ubuntu-trusty-go_agent.tgz": ExtractedPieces{
			Name:    "light-bosh",
			Version: "2579.1",

			InfName: "aws",
			HvName:  "xen-hvm",

			OSName:    "ubuntu",
			OSVersion: "trusty",

			AgentType: "go_agent",
		},

		// Notice no folder prefix
		"bosh-stemcell-2776-warden-boshlite-centos-go_agent.tgz": ExtractedPieces{
			Name:    "bosh",
			Version: "2776",

			InfName: "warden",
			HvName:  "boshlite",

			OSName:    "centos",
			OSVersion: "",

			AgentType: "go_agent",
		},
	}

	for p, e := range examples {
		path := p
		example := e

		It(fmt.Sprintf("correctly interprets '%s'", path), func() {
			s3Stemcell := NewS3Stemcell(path, "", 0, "", "")
			Expect(s3Stemcell).ToNot(BeNil())

			Expect(s3Stemcell.Name()).To(Equal(example.Name))
			Expect(s3Stemcell.Version().AsString()).To(Equal(example.Version))

			Expect(s3Stemcell.InfName()).To(Equal(example.InfName))
			Expect(s3Stemcell.HvName()).To(Equal(example.HvName))

			Expect(s3Stemcell.OSName()).To(Equal(example.OSName))
			Expect(s3Stemcell.OSVersion()).To(Equal(example.OSVersion))

			Expect(s3Stemcell.AgentType()).To(Equal(example.AgentType))
		})
	}
})
