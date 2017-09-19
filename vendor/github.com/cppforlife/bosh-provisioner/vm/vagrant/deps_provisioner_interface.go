package vagrant

type DepsProvisioner interface {
	Provision() error
	InstallRunit() error
}
