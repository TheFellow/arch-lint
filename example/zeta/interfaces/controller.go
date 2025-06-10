package interfaces

import "github.com/TheFellow/arch-lint/example/zeta/usecase"

type Controller struct {
	Service usecase.Service
}
