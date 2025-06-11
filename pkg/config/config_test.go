package config

import (
	"os"
	"testing"

	"github.com/TheFellow/arch-lint/pkg/testutil"
)

func TestLoad_SchemaValidation(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	valid := []byte("specs:\n  - name: test\n    packages:\n      include:\n        - pkg\n    rules:\n      forbid:\n        - other\n")
	good := dir + "/good.yml"
	os.WriteFile(good, valid, 0o644)

	_, err := Load(good)
	testutil.Equals(t, err, nil)

	bad := dir + "/bad.yml"
	os.WriteFile(bad, []byte("foo: bar"), 0o644)
	_, err = Load(bad)
	testutil.ErrorIf(t, err == nil, "expected error")
}
