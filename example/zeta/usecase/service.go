package usecase

import "github.com/TheFellow/arch-lint/example/zeta/infrastructure/db"

type Service struct {
	Repo db.Repository
}
