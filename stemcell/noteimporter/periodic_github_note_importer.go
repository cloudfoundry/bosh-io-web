package noteimporter

import (
	"strings"
	"time"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	bhnotesrepo "github.com/cppforlife/bosh-hub/stemcell/notesrepo"
)

type PeriodicGithubNoteImporter struct {
	p           time.Duration
	accessToken string
	stopCh      <-chan struct{}

	notesRepo bhnotesrepo.NotesRepository

	logTag string
	logger boshlog.Logger
}

func NewPeriodicGithubNoteImporter(
	p time.Duration,
	accessToken string,
	stopCh <-chan struct{},
	notesRepo bhnotesrepo.NotesRepository,
	logger boshlog.Logger,
) PeriodicGithubNoteImporter {
	return PeriodicGithubNoteImporter{
		p:           p,
		accessToken: accessToken,
		stopCh:      stopCh,

		notesRepo: notesRepo,

		logTag: "stemcell.PeriodicGithubNoteImporter",
		logger: logger,
	}
}

func (i PeriodicGithubNoteImporter) Import() error {
	i.logger.Info(i.logTag, "Starting importing stemcell notes every '%s'", i.p)

	for {
		select {
		case <-time.After(i.p):
			i.logger.Debug(i.logTag, "Looking at the stemcell notes")

			err := i.importNotes()
			if err != nil {
				i.logger.Error(i.logTag, "Failed to import notes: '%s'", err)
			}

		case <-i.stopCh:
			i.logger.Info(i.logTag, "Stopped looking at the stemcell notes")
			return nil
		}
	}
}

func (i PeriodicGithubNoteImporter) importNotes() error {
	allGhReleases, err := i.fetchAllReleasesFromGithub()
	if err != nil {
		return err
	}

	const notePrefix = "Stemcell "

	for _, ghRelease := range allGhReleases {
		matchesName := ghRelease.Name != nil && strings.HasPrefix(*ghRelease.Name, notePrefix)

		if matchesName {
			noteRec := bhnotesrepo.NoteRec{}

			// Always overwrite bosh.io release notes with GH notes;
			// covers the case when release notes are removed from GH -> remove from bosh.import
			if ghRelease.Body != nil {
				noteRec.Content = *ghRelease.Body
			}

			stemVer := strings.TrimPrefix(*ghRelease.Name, notePrefix)

			// Ideally #Save would be on the stemcell
			err = i.notesRepo.Save(stemVer, noteRec)
			if err != nil {
				return bosherr.WrapError(err, "Saving notes for stemcell version '%v'", stemVer)
			}
		}
	}

	return nil
}

func (i PeriodicGithubNoteImporter) fetchAllReleasesFromGithub() ([]github.RepositoryRelease, error) {
	i.logger.Debug(i.logTag, "Fetching github releases for stemcells")

	conf := &oauth2.Config{}

	// Authenticated access allows for 5000 reqs/hour
	client := github.NewClient(conf.Client(nil, &oauth2.Token{AccessToken: i.accessToken}))

	var allReleases []github.RepositoryRelease

	listOpts := &github.ListOptions{PerPage: 30, Page: 0}

	for {
		releases, resp, err := client.Repositories.ListReleases("cloudfoundry", "bosh", listOpts)
		if err != nil {
			return allReleases, bosherr.WrapError(err, "Listing github releases")
		}

		// Unauthenticated access can only be used up to 60 reqs/hour
		if resp.Rate.Remaining < 50 {
			waitD := resp.Rate.Reset.Sub(time.Now())

			i.logger.Debug(i.logTag, "Sleeping for '%v' until github rate-limiting resets", waitD)
			time.Sleep(waitD)
		} else {
			i.logger.Debug(i.logTag, "Left with '%d' requests for github for this hour", resp.Rate.Remaining)
		}

		allReleases = append(allReleases, releases...)
		if resp.NextPage == 0 {
			break
		}

		if len(allReleases) > 200 {
			i.logger.Debug(i.logTag, "Found '%d' stemcell releases on github", len(allReleases))
		}

		listOpts.Page = resp.NextPage
	}

	return allReleases, nil
}
