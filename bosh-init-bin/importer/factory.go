package importer

import (
	"time"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"

	bhbibrepo "github.com/cppforlife/bosh-hub/bosh-init-bin/repo"
	bhs3 "github.com/cppforlife/bosh-hub/s3"
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
	S3BoshInitBinsRepo() bhbibrepo.S3Repository
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
		repos.S3BoshInitBinsRepo(),
		logger,
	)

	return Factory{Importer: periodicImporter}
}
