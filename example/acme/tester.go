package acme

import (
	"github.com/TheFellow/go-arch-lint/example/acme/experimental"
	"github.com/TheFellow/go-arch-lint/example/acme/internal/shared"
)

func DoTest() {
	shared.Test()
	_ = experimental.NewWidget()
}
