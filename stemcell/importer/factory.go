package importer

import (
	"time"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	bhs3 "github.com/cppforlife/bosh-hub/s3"
	bhstemsrepo "github.com/cppforlife/bosh-hub/stemcell/stemsrepo"
)

type FactoryOptions struct {
	Enabled bool
	Period  time.Duration

	Buckets []BucketOptions
}

type BucketOptions struct {
	URL    string
	CDNURL string
}

type FactoryRepos interface {
	S3StemcellsRepo() bhstemsrepo.S3StemcellsRepository
}

type Factory struct {
	Importer Importer
}

func NewFactory(options FactoryOptions, repos FactoryRepos, logger boshlog.Logger) Factory {
	if !options.Enabled {
		return Factory{Importer: NewNoopImporter(logger)}
	}

	buckets := []bhs3.Bucket{}

	for _, bucketOpts := range options.Buckets {
		bucket := bhs3.NewPlainBucket(bucketOpts.URL, logger)

		if len(bucketOpts.CDNURL) > 0 {
			urlFactory := bhs3.NewDirectURLFactory(bucketOpts.CDNURL)
			bucket = bhs3.NewCDNBucket(urlFactory, bucket)
		}

		buckets = append(buckets, bucket)
	}

	periodicImporter := NewPeriodicS3BucketImporter(
		options.Period,
		make(chan struct{}),
		buckets,
		repos.S3StemcellsRepo(),
		logger,
	)

	return Factory{Importer: periodicImporter}
}
