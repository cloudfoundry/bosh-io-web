package watcher

import (
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	bhfetcher "github.com/cppforlife/bosh-hub/release/fetcher"
	bhimpsrepo "github.com/cppforlife/bosh-hub/release/importsrepo"
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
	bhwatchersrepo "github.com/cppforlife/bosh-hub/release/watchersrepo"
)

type PeriodicWatcher struct {
	p      time.Duration
	stopCh <-chan struct{}

	releasesRepo bhrelsrepo.ReleasesRepository
	watchersRepo bhwatchersrepo.WatchersRepository
	importsRepo  bhimpsrepo.ImportsRepository
	fetcher      bhfetcher.Fetcher

	logTag string
	logger boshlog.Logger
}

func NewPeriodicWatcher(
	p time.Duration,
	stopCh <-chan struct{},
	releasesRepo bhrelsrepo.ReleasesRepository,
	watchersRepo bhwatchersrepo.WatchersRepository,
	importsRepo bhimpsrepo.ImportsRepository,
	fetcher bhfetcher.Fetcher,
	logger boshlog.Logger,
) PeriodicWatcher {
	return PeriodicWatcher{
		p:      p,
		stopCh: stopCh,

		releasesRepo: releasesRepo,
		watchersRepo: watchersRepo,
		importsRepo:  importsRepo,
		fetcher:      fetcher,

		logTag: "PeriodicWatcher",
		logger: logger,
	}
}

func (w PeriodicWatcher) Watch() error {
	w.logger.Info(w.logTag, "Starting watching releases every %s", w.p)

	for {
		select {
		case <-time.After(w.p):
			err := w.lookAtReleases()
			if err != nil {
				w.logger.Error(w.logTag, "Failed to look at releases: %s", err)
			}

		case <-w.stopCh:
			w.logger.Info(w.logTag, "Stopped looking at releases")
			return nil
		}
	}
}

func (w PeriodicWatcher) lookAtReleases() error {
	w.logger.Debug(w.logTag, "Looking at releases")

	watcherRecs, err := w.watchersRepo.ListAll()
	if err != nil {
		return bosherr.WrapError(err, "Listing all watcher records")
	}

	for _, watcherRec := range watcherRecs {
		err := w.lookAtRelease(watcherRec)
		if err != nil {
			w.logger.Error(w.logTag, "Failed to look at release source '%s': %s", watcherRec.RelSource, err)
		}
	}

	return nil
}

func (w PeriodicWatcher) lookAtRelease(watcherRec bhwatchersrepo.WatcherRec) error {
	w.logger.Debug(w.logTag, "Looking at release '%v'", watcherRec)

	releaseDir, err := w.fetcher.Fetch(watcherRec.RelSource)
	if err != nil {
		return bosherr.WrapError(err, "Fetching release")
	}

	defer releaseDir.Close()

	watcherMinVersion := watcherRec.MinVersion()

	pathToManifests, err := releaseDir.ReleaseManifests()
	if err != nil {
		return bosherr.WrapError(err, "Finding release manifests")
	}

	for _, manifest := range pathToManifests {
		relVerRec, err := w.releasesRepo.Find(watcherRec.RelSource, manifest.Release.Version)
		if err != nil {
			return bosherr.WrapErrorf(err, "Finding release version '%v'", relVerRec)
		}

		if relVerRec.Version().IsLt(watcherMinVersion) {
			w.logger.Debug(w.logTag,
				"Skipping release version '%v' because it is less than watcher minimum version", relVerRec)
			continue
		}

		found, err := w.releasesRepo.Contains(relVerRec)
		if err != nil {
			return bosherr.WrapErrorf(err, "Finding release version '%v'", relVerRec)
		} else if found {
			w.logger.Debug(w.logTag,
				"Skipping release version '%v' because it is was already imported", relVerRec)
			continue
		}

		err = w.importsRepo.Push(watcherRec.RelSource, manifest.Release.Version)
		if err != nil {
			return bosherr.WrapErrorf(err, "Adding release version '%v' to import queue", relVerRec)
		}
	}

	return nil
}
