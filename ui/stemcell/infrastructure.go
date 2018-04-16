package stemcell

type Infrastructure struct {
	Name                    string
	Title                   string
	DocumentationURL        string
	LightStemcellsPublished bool

	SupportedHypervisors InfrastructureHypervisors
}

type InfrastructureHypervisor struct {
	Hypervisor Hypervisor
	Deprecated bool
}

type InfrastructureHypervisors []InfrastructureHypervisor
