package main

import (
	"bytes"
	"io"
	"os/exec"
	"strings"
	"testing"

	"github.com/TheFellow/arch-lint/pkg/testutil"
)

func TestArchLint(t *testing.T) {
	t.Parallel()
	cmd := exec.Command("go", "run", ".", "-c", "./example/rules.yml")
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	err := cmd.Run()

	testutil.ErrorIf(t, err == nil, "got %v, want %v", err, "non-nil")
	testutil.Equals(t, err.Error(), "exit status 1")
	gotOut, err := io.ReadAll(stdout)
	testutil.Equals(t, err, nil)
	testutil.Equals(t, strings.Trim(string(gotOut), "\n"), strings.Trim(wantOut, "\n"))
}

var wantOut string = `
arch-lint: [app package from api only] "example/beta/bookstore/app/books/books.go" imports "example/beta/bookstore/app/authors"
arch-lint: [app package from api or other features only] "example/epsilon/bookstore/app/books/utils/bad.go" imports "example/epsilon/bookstore/app/books"
arch-lint: [no-experimental-imports] "example/alpha/tester.go" imports "example/alpha/experimental"`
