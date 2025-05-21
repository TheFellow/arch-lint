package testutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Equals(t testing.TB, got, want any, opts ...cmp.Option) {
	t.Helper()
	if len(opts) == 0 {
		opts = []cmp.Option{cmpopts.EquateErrors()}
	}
	diff := cmp.Diff(want, got, opts...)
	if diff != "" {
		t.Fatalf("mismatch -want +got:\n%v", diff)
	}
}
