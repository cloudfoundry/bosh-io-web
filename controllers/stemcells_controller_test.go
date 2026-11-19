package controllers_test

import (
	"errors"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	mart "github.com/go-martini/martini"
	martrend "github.com/martini-contrib/render"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bosh-io/web/controllers"
	bhstemsrepo "github.com/bosh-io/web/stemcell/stemsrepo"
)

type fakeStemcellsRepository struct {
	stemcells []bhstemsrepo.Stemcell
	err       error
}

func (r fakeStemcellsRepository) FindAll(string) ([]bhstemsrepo.Stemcell, error) {
	return r.stemcells, r.err
}

var _ = Describe("stemcell API and downloads", func() {
	var errorTemplates string

	BeforeEach(func() {
		var err error
		errorTemplates, err = os.MkdirTemp("", "stemcell-controller-templates")
		Expect(err).NotTo(HaveOccurred())
		Expect(os.WriteFile(filepath.Join(errorTemplates, "error.tmpl"), []byte("{{.}}"), 0600)).To(Succeed())
	})

	AfterEach(func() {
		Expect(os.RemoveAll(errorTemplates)).To(Succeed())
	})

	request := func(path string, repo fakeStemcellsRepository) *httptest.ResponseRecorder {
		m := mart.Classic()
		m.Use(martrend.Renderer(martrend.Options{Directory: errorTemplates, Extensions: []string{".tmpl"}}))
		controller := controllers.NewStemcellsController(repo, boshlog.NewLogger(boshlog.LevelNone))
		m.Get("/d/stemcells/**", controller.Download)
		m.Get("/api/v1/stemcells/**", controller.APIV1Index)
		response := httptest.NewRecorder()
		m.ServeHTTP(response, httptest.NewRequest("GET", path, nil))
		return response
	}

	stemcell := func(key, url string) bhstemsrepo.Stemcell {
		return bhstemsrepo.NewS3Stemcell(key, "", "", "", 0, "", url)
	}

	Describe("APIV1Index", func() {
		It("returns no versions for excluded distros and lists listed distros", func() {
			xenial := stemcell("light-bosh-stemcell-621.125-aws-xen-hvm-ubuntu-xenial-go_agent.tgz", "https://example.com/xenial.tgz")
			jammy := stemcell("light-bosh-stemcell-1.900-aws-xen-hvm-ubuntu-jammy-go_agent.tgz", "https://example.com/jammy.tgz")

			excluded := request("/api/v1/stemcells/bosh-aws-xen-hvm-ubuntu-xenial-go_agent", fakeStemcellsRepository{stemcells: []bhstemsrepo.Stemcell{xenial}})
			Expect(excluded.Code).To(Equal(200))
			Expect(strings.TrimSpace(excluded.Body.String())).To(Equal("[]"))

			unknown := request("/api/v1/stemcells/bosh-aws-xen-hvm-ubuntu-unknown-go_agent", fakeStemcellsRepository{})
			Expect(unknown.Code).To(Equal(200))
			Expect(strings.TrimSpace(unknown.Body.String())).To(Equal("[]"))

			listed := request("/api/v1/stemcells/bosh-aws-xen-hvm-ubuntu-jammy-go_agent", fakeStemcellsRepository{stemcells: []bhstemsrepo.Stemcell{jammy}})
			Expect(listed.Code).To(Equal(200))
			Expect(listed.Body.String()).To(ContainSubstring(`"version":"1.900"`))
		})
	})

	Describe("Download", func() {
		It("hides excluded distros and redirects listed distros", func() {
			xenial := stemcell("light-bosh-stemcell-621.125-aws-xen-hvm-ubuntu-xenial-go_agent.tgz", "https://example.com/xenial.tgz")
			jammy := stemcell("light-bosh-stemcell-1.900-aws-xen-hvm-ubuntu-jammy-go_agent.tgz", "https://example.com/jammy.tgz")

			latestExcluded := request("/d/stemcells/bosh-aws-xen-hvm-ubuntu-xenial-go_agent", fakeStemcellsRepository{stemcells: []bhstemsrepo.Stemcell{xenial}})
			Expect(latestExcluded.Code).To(Equal(404))
			Expect(strings.TrimSpace(latestExcluded.Body.String())).To(Equal("Latest stemcell is not found"))

			versionExcluded := request("/d/stemcells/bosh-aws-xen-hvm-ubuntu-xenial-go_agent?v=621.125", fakeStemcellsRepository{stemcells: []bhstemsrepo.Stemcell{xenial}})
			Expect(versionExcluded.Code).To(Equal(404))
			Expect(strings.TrimSpace(versionExcluded.Body.String())).To(Equal("Stemcell version &#39;621.125&#39; is not found"))

			latestListed := request("/d/stemcells/bosh-aws-xen-hvm-ubuntu-jammy-go_agent", fakeStemcellsRepository{stemcells: []bhstemsrepo.Stemcell{jammy}})
			Expect(latestListed.Code).To(Equal(302))
			Expect(latestListed.Header().Get("Location")).To(Equal("https://example.com/jammy.tgz"))
		})

		It("preserves invalid-version and repository-error statuses", func() {
			jammy := stemcell("light-bosh-stemcell-1.900-aws-xen-hvm-ubuntu-jammy-go_agent.tgz", "https://example.com/jammy.tgz")
			invalidVersion := request("/d/stemcells/bosh-aws-xen-hvm-ubuntu-jammy-go_agent?v=not-a-version!", fakeStemcellsRepository{stemcells: []bhstemsrepo.Stemcell{jammy}})
			Expect(invalidVersion.Code).To(Equal(400))

			repositoryError := request("/d/stemcells/bosh-aws-xen-hvm-ubuntu-jammy-go_agent", fakeStemcellsRepository{err: errors.New("repository failed")})
			Expect(repositoryError.Code).To(Equal(500))
			body, err := io.ReadAll(repositoryError.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(ContainSubstring("repository failed"))
		})
	})
})
