package templatesrepo

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	bpdep "github.com/cppforlife/bosh-provisioner/deployment"
	bpindex "github.com/cppforlife/bosh-provisioner/index"
)

type CTRepository struct {
	index  bpindex.Index
	logger boshlog.Logger
}

// todo fingerprint property changes
type jobToTemplateKey struct {
	JobName string
}

func NewConcreteTemplatesRepository(
	index bpindex.Index,
	logger boshlog.Logger,
) CTRepository {
	return CTRepository{index: index, logger: logger}
}

func (tr CTRepository) Find(job bpdep.Job, instance bpdep.Instance) (TemplateRecord, bool, error) {
	var record TemplateRecord

	err := tr.index.Find(tr.templateKey(job), &record)
	if err != nil {
		if err == bpindex.ErrNotFound {
			return record, false, nil
		}

		return record, false, bosherr.WrapError(err, "Finding tempate")
	}

	return record, true, nil
}

func (tr CTRepository) Save(job bpdep.Job, instance bpdep.Instance, record TemplateRecord) error {
	err := tr.index.Save(tr.templateKey(job), record)
	if err != nil {
		return bosherr.WrapError(err, "Saving template")
	}

	return nil
}

func (tr CTRepository) templateKey(job bpdep.Job) jobToTemplateKey {
	return jobToTemplateKey{JobName: job.Name}
}
