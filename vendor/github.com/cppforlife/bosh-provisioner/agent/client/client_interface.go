package client

import (
	boshaction "github.com/cloudfoundry/bosh-agent/agent/action"
	boshas "github.com/cloudfoundry/bosh-agent/agent/applier/applyspec"
	boshcomp "github.com/cloudfoundry/bosh-agent/agent/compiler"
)

type Client interface {
	TaskManager
	VMAdministrator
	StateManager
	JobManager
	PackageCompiler
	DiskManager
	NetworkManager
}

type TaskManager interface {
	Ping() (string, error)
	GetTask(string) (interface{}, error)
	CancelTask(string) (string, error)
}

// VMAdministrator provides administrative API for the agent
// todo eliminate remaining param names
type VMAdministrator interface {
	SSH(cmd string, params boshaction.SSHParams) (map[string]interface{}, error)
	FetchLogs(logType string, filters []string) (map[string]interface{}, error)
}

type StateManager interface {
	Prepare(boshas.V1ApplySpec) (string, error)
	Apply(boshas.V1ApplySpec) (string, error)
	GetState(filters ...string) (boshaction.GetStateV1ApplySpec, error)
}

type JobManager interface {
	PreStart() error
	Start() (string, error)
	PostStart() error
	Stop() (string, error)
	Drain(boshaction.DrainType, ...boshas.V1ApplySpec) (int, error)
	RunErrand() (boshaction.ErrandResult, error)
}

// CompiledPackage keeps information about compiled package asset
// todo move into agent
type CompiledPackage struct {
	BlobID string `json:"blobstore_id"`
	SHA1   string `json:"sha1"`
}

type PackageCompiler interface {
	CompilePackage(
		blobID string,
		sha1 string,
		name string,
		version string,
		deps boshcomp.Dependencies,
	) (CompiledPackage, error)
}

type DiskManager interface {
	// ListDisk(...)
	// MigrateDisk(...)
	// MountDisk(...)
	// UnmountDisk(...)
}

type NetworkManager interface {
	// PrepareNetworkChange(...)
	// PrepareConfigureNetworks(...)
	// ConfigureNetworks(...)
}
