package alpha

import (
	"github.com/TheFellow/arch-lint/example/alpha/experimental"
	"github.com/TheFellow/arch-lint/example/alpha/internal/exception"
)

func DoTest() {
	exception.Test()
	_ = experimental.NewWidget()
}
