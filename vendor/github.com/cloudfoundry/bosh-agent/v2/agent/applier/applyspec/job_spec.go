package applyspec

import (
	models "github.com/cloudfoundry/bosh-agent/v2/agent/applier/models"
)

type JobSpec struct {
	Name             *string           `json:"name"`
	Release          string            `json:"release"`
	Template         string            `json:"template"`
	Version          string            `json:"version"`
	JobTemplateSpecs []JobTemplateSpec `json:"templates"`
}

func (s *JobSpec) JobTemplateSpecsAsJobs() []models.Job {
	jobs := []models.Job{}
	for _, value := range s.JobTemplateSpecs {
		jobs = append(jobs, value.AsJob())
	}
	return jobs
}
