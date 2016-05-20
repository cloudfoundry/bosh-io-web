package noteimporter

import (
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	bhnotesrepo "github.com/cppforlife/bosh-hub/release/notesrepo"
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

type PeriodicGithubNoteImporter struct {
	p           time.Duration
	accessToken string
	stopCh      <-chan struct{}

	releasesRepo bhrelsrepo.ReleasesRepository

	logTag string
	logger boshlog.Logger
}

func NewPeriodicGithubNoteImporter(
	p time.Duration,
	accessToken string,
	stopCh <-chan struct{},
	releasesRepo bhrelsrepo.ReleasesRepository,
	logger boshlog.Logger,
) PeriodicGithubNoteImporter {
	return PeriodicGithubNoteImporter{
		p:           p,
		accessToken: accessToken,
		stopCh:      stopCh,

		releasesRepo: releasesRepo,

		logTag: "PeriodicGithubNoteImporter",
		logger: logger,
	}
}

func (i PeriodicGithubNoteImporter) Import() error {
	i.logger.Info(i.logTag, "Starting importing release notes every '%s'", i.p)

	for {
		select {
		case <-time.After(i.p):
			i.logger.Debug(i.logTag, "Looking at the release versions")

			err := i.importNotes()
			if err != nil {
				i.logger.Error(i.logTag, "Failed to import notes: '%s'", err)
			}

		case <-i.stopCh:
			i.logger.Info(i.logTag, "Stopped looking at the release versions")
			return nil
		}
	}
}

func (i PeriodicGithubNoteImporter) importNotes() error {
	// todo getting all releases; could be big
	sources, err := i.releasesRepo.ListAll()
	if err != nil {
		return bosherr.WrapError(err, "Listing releases")
	}

	for _, source := range sources {
		ghSource, valid := newGithubSource(source)
		if !valid {
			continue
		}

		// todo convert Source to string; argh
		relVerRecs, err := i.releasesRepo.FindAll(source.Full)
		if err != nil {
			return bosherr.WrapError(err, "Listing all versions for release source '%s'", source)
		}

		err = i.importNotesForRelease(ghSource, relVerRecs)
		if err != nil {
			return bosherr.WrapError(err, "Importing notes for release source '%s'", source)
		}
	}

	return nil
}

func (i PeriodicGithubNoteImporter) importNotesForRelease(ghSource githubSource, relVerRecs []bhrelsrepo.ReleaseVersionRec) error {
	// fast path if there are no release versions
	if len(relVerRecs) == 0 {
		return nil
	}

	allGhReleases, err := i.fetchAllReleasesFromGithub(ghSource)
	if err != nil {
		// Continue onto other release versions
		i.logger.Error(i.logTag, "Failed to fetch releases from github for '%v': %s", ghSource, err)
		return nil
	}

	for _, relVerRec := range relVerRecs {
		for _, ghRelease := range allGhReleases {
			expectedLabel := "v" + relVerRec.VersionRaw

			// Either release name or git tag name match
			matchesName := ghRelease.Name != nil && *ghRelease.Name == expectedLabel
			matchesTagName := ghRelease.TagName != nil && *ghRelease.TagName == expectedLabel

			if matchesName || matchesTagName {
				noteRec := bhnotesrepo.NoteRec{}

				// Always overwrite bosh.io release notes with GH notes;
				// covers the case when release notes are removed from GH -> remove from bosh.import
				if ghRelease.Body != nil {
					noteRec.Content = *ghRelease.Body
				}

				err = relVerRec.SetNotes(noteRec)
				if err != nil {
					return bosherr.WrapError(err, "Saving notes for release version '%v'", relVerRec)
				}

				break
			}
		}
	}

	return nil
}

func (i PeriodicGithubNoteImporter) fetchAllReleasesFromGithub(ghSource githubSource) ([]github.RepositoryRelease, error) {
	i.logger.Debug(i.logTag, "Fetching github releases for '%v'", ghSource)

	conf := &oauth2.Config{}

	// Authenticated access allows for 5000 reqs/hour
	client := github.NewClient(conf.Client(nil, &oauth2.Token{AccessToken: i.accessToken}))

	var allReleases []github.RepositoryRelease

	listOpts := &github.ListOptions{PerPage: 30, Page: 0}

	for {
		releases, resp, err := client.Repositories.ListReleases(ghSource.Owner, ghSource.Repo, listOpts)
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
			i.logger.Debug(i.logTag, "Found '%d' releases on github for '%v'", len(allReleases), ghSource)
		}

		listOpts.Page = resp.NextPage
	}

	return allReleases, nil
}
