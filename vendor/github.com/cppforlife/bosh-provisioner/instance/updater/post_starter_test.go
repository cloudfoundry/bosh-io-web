package updater_test

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	fakebpagclient "github.com/cppforlife/bosh-provisioner/agent/client/fakes"
	. "github.com/cppforlife/bosh-provisioner/instance/updater"
)

var _ = Describe("PostStarter", func() {
	var (
		agentClient *fakebpagclient.FakeClient
		logger      boshlog.Logger
		postStarter PostStarter
	)

	BeforeEach(func() {
		agentClient = &fakebpagclient.FakeClient{}
		logger = boshlog.NewLogger(boshlog.LevelNone)
		postStarter = NewPostStarter(agentClient, logger)
	})

	Describe("PostStart", func() {
		Context("when the script exits with a 0 exit status", func() {
			It("does not return an error", func() {
				err := postStarter.PostStart()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when the script exits with a non 0 exit status", func() {
			BeforeEach(func() {
				agentClient.PostStartErr = bosherr.Error("some-error")
			})

			It("return an error", func() {
				err := postStarter.PostStart()
				Expect(err).To(MatchError("Post-Starting: some-error"))
			})
		})
	})
})
