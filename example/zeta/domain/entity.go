package domain

import "github.com/TheFellow/arch-lint/example/zeta/util"

type Entity struct {
}

func (e Entity) Validate() (bool, error) {
	util.SomeHelper() // domain should have no dependencies
	return true, nil
}
