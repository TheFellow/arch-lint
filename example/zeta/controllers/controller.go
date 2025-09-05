package controllers

import (
	"github.com/TheFellow/arch-lint/example/zeta/infrastructure/db"
)

type Controller struct {
	// Service usecase.Service // should use this
	Repo db.Repository // and not this
}
