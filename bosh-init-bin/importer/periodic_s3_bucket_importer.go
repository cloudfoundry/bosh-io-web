package importer

import (
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	bhbibrepo "github.com/cppforlife/bosh-hub/bosh-init-bin/repo"
	bhs3 "github.com/cppforlife/bosh-hub/s3"
)

type PeriodicS3BucketImporter struct {
	p      time.Duration
	stopCh <-chan struct{}

	buckets          []bhs3.Bucket
	boshInitBinsRepo bhbibrepo.S3Repository

	logTag string
	logger boshlog.Logger
}

func NewPeriodicS3BucketImporter(
	p time.Duration,
	stopCh <-chan struct{},
	buckets []bhs3.Bucket,
	boshInitBinsRepo bhbibrepo.S3Repository,
	logger boshlog.Logger,
) PeriodicS3BucketImporter {
	return PeriodicS3BucketImporter{
		p:      p,
		stopCh: stopCh,

		buckets:          buckets,
		boshInitBinsRepo: boshInitBinsRepo,

		logTag: "PeriodicS3BucketImporter",
		logger: logger,
	}
}

func (i PeriodicS3BucketImporter) Import() error {
	i.logger.Info(i.logTag, "Starting importing S3 bucket bosh-init bins every '%s'", i.p)

	for {
		select {
		case <-time.After(i.p):
			i.logger.Debug(i.logTag, "Looking at the bucket")

			err := i.importAll()
			if err != nil {
				i.logger.Error(i.logTag, "Failed to import bucket: '%s'", err)
			}

		case <-i.stopCh:
			i.logger.Info(i.logTag, "Stopped looking at the bucket")
			return nil
		}
	}
}

func (i PeriodicS3BucketImporter) importAll() error {
	allFiles := []bhs3.File{}

	for _, bucket := range i.buckets {
		someFiles, err := bucket.Files()
		if err != nil {
			// Return immediately if it cannot fetch latest files
			// from any bucket so that saved binaries do not get overwritten
			return bosherr.WrapError(err, "Getting bucket files")
		}

		allFiles = append(allFiles, someFiles...)
	}

	err := i.boshInitBinsRepo.SaveAll(allFiles)
	if err != nil {
		return bosherr.WrapError(err, "Saving bucket files in bosh-init bins repo")
	}

	return nil
}
