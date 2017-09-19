package templatesrepo

import (
	bpdep "github.com/cppforlife/bosh-provisioner/deployment"
)

type TemplateRecord struct {
	BlobID string
	SHA1   string
}

// TemplatesRepository maintains list of rendered templates as blobs
type TemplatesRepository interface {
	Find(bpdep.Job, bpdep.Instance) (TemplateRecord, bool, error)
	Save(bpdep.Job, bpdep.Instance, TemplateRecord) error
}
