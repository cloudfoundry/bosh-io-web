package stemcell

import "fmt"

type Infrastructures []Infrastructure

func (i Infrastructures) ByName(name string) (Infrastructure, error) {
	for _, infrastructure := range i {
		if infrastructure.Name == name {
			return infrastructure, nil
		}
	}

	return Infrastructure{}, fmt.Errorf("unknown infrastructure: %s", name)
}
