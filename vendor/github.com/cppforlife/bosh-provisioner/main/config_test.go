package main_test

import (
	fakesys "github.com/cloudfoundry/bosh-utils/system/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cppforlife/bosh-provisioner/main"
	bpvm "github.com/cppforlife/bosh-provisioner/vm"
)

var _ = Describe("NewConfigFromPath", func() {
	var (
		fs *fakesys.FakeFileSystem
	)

	BeforeEach(func() {
		fs = fakesys.NewFakeFileSystem()
	})

	It("defautls null values to agent provisioner config defaults", func() {
		configJSON := `{
      "manifest_path": "fake-manifest-path",
      "assets_dir": "fake-assets-dir",
      "repos_dir": "fake-repos-dir",
      "blobstore": {
        "provider": "local",
        "options": {
          "blobstore_path": "fake-blobstore-path"
        }
      },
      "vm_provisioner": {
        "agent_provisioner": {
          "platform": null,
          "configuration": null,
          "mbus": null
        }
      }
    }`

		err := fs.WriteFileString("/tmp/config", configJSON)
		Expect(err).ToNot(HaveOccurred())

		config, err := NewConfigFromPath("/tmp/config", fs)
		Expect(err).ToNot(HaveOccurred())

		Expect(config.VMProvisioner.AgentProvisioner).To(Equal(
			bpvm.AgentProvisionerConfig{
				Platform: "ubuntu",

				Configuration: map[string]interface{}{
					"Infrastructure": map[string]interface{}{
						"Settings": map[string]interface{}{
							"UseRegistry": true,
							"Sources": []map[string]interface{}{
								{
									"SettingsPath": "warden-cpi-agent-env.json",
									"Type":         "File",
								},
							},
						},
					},
					"Platform": map[string]interface{}{
						"Linux": map[string]interface{}{
							"UseDefaultTmpDir": true,
							"SkipDiskSetup":    true,
						},
					},
				},

				Mbus: "https://user:password@127.0.0.1:4321/agent",
			},
		))
	})
})
