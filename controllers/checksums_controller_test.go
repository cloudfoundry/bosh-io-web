package controllers_test

import (
	"net/http"
	"net/url"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	mart "github.com/go-martini/martini"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	bhchecks "github.com/cppforlife/bosh-hub/checksumsrepo"
	. "github.com/cppforlife/bosh-hub/controllers"
)

var _ = Describe("ChecksumsController", func() {
	var (
		allowedMatches []ChecksumReqMatch
		repo           *FakeChecksumsRepository
		controller     ChecksumsController

		renderer *FakeRender
	)

	BeforeEach(func() {
		allowedMatches = make([]ChecksumReqMatch, 100)
		repo = &FakeChecksumsRepository{}
		logger := boshlog.NewLogger(boshlog.LevelNone)
		controller = NewChecksumsController(allowedMatches, repo, logger)

		renderer = &FakeRender{}
	})

	Describe("Save", func() {
		It("returns 401 if Authorization header is missing", func() {
			controller.Save(&http.Request{}, renderer, mart.Params{"_1": "matches-key"})
			Expect(renderer.JSONStatus).To(Equal(401))
			Expect(renderer.JSONResponse).To(Equal(map[string]string{"error": "Unauthorized: Token must be non-empty"}))
		})

		It("returns 401 if Authorization header is empty", func() {
			req := &http.Request{
				Header: http.Header{},
			}
			req.Header.Add("Authorization", "")

			controller.Save(req, renderer, mart.Params{"_1": "matches-key"})
			Expect(renderer.JSONStatus).To(Equal(401))
			Expect(renderer.JSONResponse).To(Equal(map[string]string{"error": "Unauthorized: Token must be non-empty"}))
		})

		It("returns 401 if Authorization header doesnt match token for matching key", func() {
			allowedMatches[0] = ChecksumReqMatch{
				Token:    "allowed-token",
				Includes: []string{"key"},
			}

			req := &http.Request{
				Header: http.Header{},
			}
			req.Header.Add("Authorization", "bearer given-token")

			controller.Save(req, renderer, mart.Params{"_1": "matches-key"})
			Expect(renderer.JSONStatus).To(Equal(401))
			Expect(renderer.JSONResponse).To(Equal(map[string]string{"error": "Unauthorized: Token mismatch"}))
		})

		It("returns 401 if Authorization header matches token but also matches more than one key", func() {
			allowedMatches[0] = ChecksumReqMatch{
				Token:    "allowed-token",
				Includes: []string{"key"},
			}

			allowedMatches[1] = ChecksumReqMatch{
				Token:    "other-allowed-token",
				Includes: []string{"key"},
			}

			validSHA1 := "sha1-for-matches-key"

			req := &http.Request{
				Header: http.Header{},
				Form:   url.Values{"sha1": []string{validSHA1}},
			}
			req.Header.Add("Authorization", "bearer allowed-token")

			renderer = &FakeRender{}
			controller.Save(req, renderer, mart.Params{"_1": "matches-key"})
			Expect(renderer.JSONStatus).To(Equal(200))
			Expect(renderer.JSONResponse).To(Equal(map[string]string{"sha1": validSHA1}))

			req = &http.Request{
				Header: http.Header{},
				Form:   url.Values{"sha1": []string{validSHA1}},
			}
			req.Header.Add("Authorization", "bearer other-allowed-token")

			renderer = &FakeRender{}
			controller.Save(req, renderer, mart.Params{"_1": "matches-key"})
			Expect(renderer.JSONStatus).To(Equal(200))
			Expect(renderer.JSONResponse).To(Equal(map[string]string{"sha1": validSHA1}))
		})

		It("returns 401 if Authorization header matches token without matching key", func() {
			allowedMatches[0] = ChecksumReqMatch{
				Token:    "allowed-token",
				Includes: []string{"key"},
			}

			req := &http.Request{
				Header: http.Header{},
			}
			req.Header.Add("Authorization", "bearer allowed-token")

			controller.Save(req, renderer, mart.Params{"_1": "does-not-match-KEY"}) // because it's uppercase
			Expect(renderer.JSONStatus).To(Equal(401))
			Expect(renderer.JSONResponse).To(Equal(map[string]string{"error": "Unauthorized: Token mismatch"}))
		})

		It("returns 401 if key is empty", func() {
			allowedMatches[0] = ChecksumReqMatch{
				Token:    "allowed-token",
				Includes: []string{"key"},
			}

			req := &http.Request{
				Header: http.Header{},
			}
			req.Header.Add("Authorization", "bearer allowed-token")

			controller.Save(req, renderer, mart.Params{"_1": ""})
			Expect(renderer.JSONStatus).To(Equal(401))
			Expect(renderer.JSONResponse).To(Equal(map[string]string{"error": "Unauthorized: Key must be non-empty"}))
		})

		It("return 200 and saves checksum if token and key matches", func() {
			allowedMatches[0] = ChecksumReqMatch{
				Token:    "allowed-token",
				Includes: []string{"key"},
			}

			validSHA1 := "sha1-for-matches-key"

			req := &http.Request{
				Header: http.Header{},
				Form:   url.Values{"sha1": []string{validSHA1}},
			}
			req.Header.Add("Authorization", "bearer allowed-token")

			controller.Save(req, renderer, mart.Params{"_1": "matches-key"})
			Expect(renderer.JSONStatus).To(Equal(200))
			Expect(renderer.JSONResponse).To(Equal(map[string]string{"sha1": validSHA1}))

			Expect(repo.SavedKey).To(Equal("matches-key"))
			Expect(repo.SavedRecord).To(Equal(bhchecks.ChecksumRec{SHA1: validSHA1}))
		})
	})
})
