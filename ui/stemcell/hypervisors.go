package stemcell

import "fmt"

type Hypervisors []Hypervisor

func (i Hypervisors) ByName(name string) (Hypervisor, error) {
	for _, hypervisor := range i {
		if hypervisor.Name == name {
			return hypervisor, nil
		}
	}

	return Hypervisor{}, fmt.Errorf("unknown hypervisor: %s", name)
}
