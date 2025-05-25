package linter

import (
	"testing"

	"github.com/TheFellow/arch-lint/pkg/testutil"
)

func TestViolation_String(t *testing.T) {
	t.Parallel()
	v := Violation{
		Rule:   "my-rule",
		File:   "foo.go",
		Import: "path/to/bar",
	}
	got := v.String()
	want := `arch-lint: [my-rule] "foo.go" imports "path/to/bar"`
	testutil.Equals(t, got, want)
}
