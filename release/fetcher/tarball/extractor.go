package tarball

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	bprel "github.com/cppforlife/bosh-provisioner/release"
	bpreljob "github.com/cppforlife/bosh-provisioner/release/job"

	bhjobsrepo "github.com/cppforlife/bosh-hub/release/jobsrepo"
	bhrelsrepo "github.com/cppforlife/bosh-hub/release/releasesrepo"
)

type Extractor struct {
	releaseReaderFactory bprel.ReaderFactory
	jobReaderFactory     bpreljob.ReaderFactory

	releasesRepo        bhrelsrepo.ReleasesRepository
	releaseVersionsRepo bhrelsrepo.ReleaseVersionsRepository
	jobsRepo            bhjobsrepo.JobsRepository

	logTag string
	logger boshlog.Logger
}

func NewExtractor(
	releaseReaderFactory bprel.ReaderFactory,
	jobReaderFactory bpreljob.ReaderFactory,
	releasesRepo bhrelsrepo.ReleasesRepository,
	releaseVersionsRepo bhrelsrepo.ReleaseVersionsRepository,
	jobsRepo bhjobsrepo.JobsRepository,
	logger boshlog.Logger,
) Extractor {
	return Extractor{
		releaseReaderFactory: releaseReaderFactory,
		jobReaderFactory:     jobReaderFactory,

		releasesRepo:        releasesRepo,
		releaseVersionsRepo: releaseVersionsRepo,
		jobsRepo:            jobsRepo,

		logTag: "Extractor",
		logger: logger,
	}
}

func (e Extractor) Extract(url, tgzPath string) (bhrelsrepo.ReleaseVersionRec, error) {
	e.logger.Debug(e.logTag, "Extracting from '%s' (url '%s')", tgzPath, url)

	var relVerRec bhrelsrepo.ReleaseVersionRec

	rel, relJobs, err := e.extractReleaseAndJobs(tgzPath)
	if err != nil {
		return relVerRec, bosherr.WrapError(err, "Extracting release and jobs")
	}

	relVerRec, err = e.releasesRepo.Find(url, rel.Version)
	if err != nil {
		return relVerRec, bosherr.WrapError(err, "Finding release")
	}

	err = e.jobsRepo.SaveAll(relVerRec, relJobs)
	if err != nil {
		return relVerRec, bosherr.WrapError(err, "Saving release jobs into jobs repository")
	}

	err = e.releaseVersionsRepo.Save(relVerRec, rel)
	if err != nil {
		return relVerRec, bosherr.WrapError(err, "Saving release into releases repository")
	}

	return relVerRec, nil
}

func (e Extractor) extractReleaseAndJobs(tgzPath string) (bprel.Release, []bpreljob.Job, error) {
	var rel bprel.Release

	relReader := e.releaseReaderFactory.NewTarReader("file://" + tgzPath)

	rel, err := relReader.Read()
	if err != nil {
		return rel, nil, bosherr.WrapError(err, "Reading release")
	}

	defer relReader.Close()

	var relJobs []bpreljob.Job

	for _, j := range rel.Jobs {
		relJobReader := e.jobReaderFactory.NewTarReader("file://" + j.TarPath)

		relJob, err := relJobReader.Read()
		if err != nil {
			return rel, nil, bosherr.WrapErrorf(err, "Reading release job '%s'", j.Name)
		}

		defer relJobReader.Close()

		relJobs = append(relJobs, relJob)
	}

	return rel, relJobs, nil
}
