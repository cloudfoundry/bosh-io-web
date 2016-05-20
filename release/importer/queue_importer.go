package importer

import (
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	bhfetcher "github.com/cppforlife/bosh-hub/release/fetcher"
	bhimperrsrepo "github.com/cppforlife/bosh-hub/release/importerrsrepo"
	bhimpsrepo "github.com/cppforlife/bosh-hub/release/importsrepo"
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

type QueueImporter struct {
	p      time.Duration
	stopCh <-chan struct{}

	tarballReleaseFactory bhfetcher.TarballReleaseFactory
	releasesRepo          bhrelsrepo.ReleasesRepository
	importsRepo           bhimpsrepo.ImportsRepository
	importErrsRepo        bhimperrsrepo.ImportErrsRepository
	fetcher               bhfetcher.Fetcher

	logTag string
	logger boshlog.Logger
}

func NewQueueImporter(
	p time.Duration,
	stopCh <-chan struct{},
	tarballReleaseFactory bhfetcher.TarballReleaseFactory,
	releasesRepo bhrelsrepo.ReleasesRepository,
	importsRepo bhimpsrepo.ImportsRepository,
	importErrsRepo bhimperrsrepo.ImportErrsRepository,
	fetcher bhfetcher.Fetcher,
	logger boshlog.Logger,
) QueueImporter {
	return QueueImporter{
		p:      p,
		stopCh: stopCh,

		tarballReleaseFactory: tarballReleaseFactory,
		releasesRepo:          releasesRepo,
		importsRepo:           importsRepo,
		importErrsRepo:        importErrsRepo,
		fetcher:               fetcher,

		logTag: "QueueImporter",
		logger: logger,
	}
}

func (i QueueImporter) Import() error {
	i.logger.Info(i.logTag, "Starting importing releases every '%s'", i.p)

	for {
		select {
		case <-time.After(i.p):
			i.logger.Debug(i.logTag, "Looking at the queue")

			for {
				stop, err := i.pickUpNext()
				if err != nil {
					i.logger.Error(i.logTag, "Failed to pick up next: '%s'", err)
				}

				if stop {
					break
				}
			}

		case <-i.stopCh:
			i.logger.Info(i.logTag, "Stopped looking at the queue")
			return nil
		}
	}
}

func (i QueueImporter) pickUpNext() (bool, error) {
	importRec, found, err := i.importsRepo.Pull()
	if err != nil {
		return true, bosherr.WrapError(err, "Listing all watcher records")
	} else if !found {
		return true, nil // There are no more pending imports
	}

	// todo releasesRepo.Find vs Contains is weird
	relVerRec, err := i.releasesRepo.Find(importRec.RelSource, importRec.Version)
	if err != nil {
		return false, bosherr.WrapError(err, "Finding release version '%v'", relVerRec)
	}

	found, err = i.releasesRepo.Contains(relVerRec)
	if err != nil {
		return false, bosherr.WrapError(err, "Contains check release version '%v'", relVerRec)
	} else if found {
		return false, nil // Already imported by someone else; skipping
	}

	found, err = i.importErrsRepo.Contains(importRec)
	if err != nil {
		return false, bosherr.WrapError(err, "Finding import err '%v'", importRec)
	} else if found {
		return false, nil // Cannot import until resolved
	}

	return false, i.importRelease(importRec)
}

func (i QueueImporter) importRelease(importRec bhimpsrepo.ImportRec) error {
	i.logger.Debug(i.logTag, "Planning to import release '%v'", importRec)

	releaseDir, err := i.fetcher.Fetch(importRec.RelSource)
	if err != nil {
		return bosherr.WrapError(err, "Fetching release source '%s'", importRec.RelSource)
	}

	defer releaseDir.Close()

	pathToManifests, err := releaseDir.ReleaseManifests()
	if err != nil {
		return bosherr.WrapError(err, "Finding release manifests")
	}

	for manifestPath, manifest := range pathToManifests {
		// todo only filtering by version; should be name+version
		if manifest.Release.Version != importRec.Version {
			continue
		}

		tarballRel := i.tarballReleaseFactory.NewTarballRelease(manifestPath)

		trackStopCh := make(chan struct{})

		go i.trackProgress(importRec, trackStopCh)

		err = tarballRel.Import(importRec.RelSource)
		if err != nil {
			importErr := bhimperrsrepo.ImportErrRec{
				ImportRec: importRec,
				Err:       err.Error(),
			}

			saveErr := i.importErrsRepo.Add(importErr)
			if saveErr != nil {
				i.logger.Debug(i.logTag, "Failed to add import err '%s' for import '%v'", importRec, saveErr.Error())
			}

			trackStopCh <- struct{}{}
			return bosherr.WrapError(err, "Importing tarball release")
		}

		trackStopCh <- struct{}{}

		i.logger.Debug(i.logTag, "Successfully imported release '%v'", importRec)

		return nil
	}

	i.logger.Debug(i.logTag, "Did not find release '%v' to import", importRec)

	return nil
}

func (i QueueImporter) trackProgress(importRec bhimpsrepo.ImportRec, stopCh <-chan struct{}) {
	for j := 0; ; j++ {
		select {
		case <-time.After(1 * time.Second):
			i.logger.Debug(i.logTag, "Tracking progress of '%v': step '%d'", importRec, j)

		case <-stopCh:
			return
		}
	}
}
