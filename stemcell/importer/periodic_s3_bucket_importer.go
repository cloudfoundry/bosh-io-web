package importer

import (
	"time"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	bhs3 "github.com/cppforlife/bosh-hub/s3"
	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
)

type PeriodicS3BucketImporter struct {
	p      time.Duration
	stopCh <-chan struct{}

	buckets       []bhs3.Bucket
	stemcellsRepo bhstemsrepo.S3StemcellsRepository

	logTag string
	logger boshlog.Logger
}

func NewPeriodicS3BucketImporter(
	p time.Duration,
	stopCh <-chan struct{},
	buckets []bhs3.Bucket,
	stemcellsRepo bhstemsrepo.S3StemcellsRepository,
	logger boshlog.Logger,
) PeriodicS3BucketImporter {
	return PeriodicS3BucketImporter{
		p:      p,
		stopCh: stopCh,

		buckets:       buckets,
		stemcellsRepo: stemcellsRepo,

		logTag: "PeriodicS3BucketImporter",
		logger: logger,
	}
}

func (i PeriodicS3BucketImporter) Import() error {
	i.logger.Info(i.logTag, "Starting importing S3 bucket stemcells every '%s'", i.p)

	for {
		select {
		case <-time.After(i.p):
			i.logger.Debug(i.logTag, "Looking at the bucket")

			err := i.importAllStemcells()
			if err != nil {
				i.logger.Error(i.logTag, "Failed to import bucket: '%s'", err)
			}

		case <-i.stopCh:
			i.logger.Info(i.logTag, "Stopped looking at the bucket")
			return nil
		}
	}
}

func (i PeriodicS3BucketImporter) importAllStemcells() error {
	allFiles := []bhs3.File{}

	for _, bucket := range i.buckets {
		someFiles, err := bucket.Files()
		if err != nil {
			// Return immediately if it cannot fetch latest files
			// from any bucket so that saved stemcells do not get overwritten
			return bosherr.WrapError(err, "Getting bucket files")
		}

		allFiles = append(allFiles, someFiles...)
	}

	err := i.stemcellsRepo.SaveAll(allFiles)
	if err != nil {
		return bosherr.WrapError(err, "Saving bucket files in stemcells repo")
	}

	return nil
}
