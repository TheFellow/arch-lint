package domain

import "github.com/TheFellow/arch-lint/example/zeta/usecase"

type Entity struct {
}

var _ = usecase.Service{}

func (e Entity) Validate() (bool, error) {
	return true, nil
}
