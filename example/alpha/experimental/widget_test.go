package experimental_test

import (
	"testing"

	"github.com/TheFellow/arch-lint/example/alpha/experimental"
	"github.com/TheFellow/arch-lint/pkg/testutil"
)

func TestNewWidget(t *testing.T) {
	w := experimental.NewWidget()
	testutil.ErrorIf(t, w == nil, "got %v, want %v", w, "non-nil")
}
