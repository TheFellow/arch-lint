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
	got := strings.Split(strings.Trim(string(gotOut), "\n"), "\n")
	want := strings.Split(strings.Trim(wantOut, "\n"), "\n")
	testutil.Equals(t, got, want)
}

var wantOut string = `
arch-lint: [app package from api only] package "example/beta/bookstore/app/books" imports "example/beta/bookstore/app/authors"
arch-lint: [app package from api or other features only] package "example/epsilon/bookstore/app/books/utils" imports "example/epsilon/bookstore/app/books"
arch-lint: [clean architecture - controllers without infrastructure] package "example/zeta/controllers" imports "example/zeta/infrastructure/db"
arch-lint: [clean architecture - domain independent] package "example/zeta/domain" imports "example/zeta/usecase"
arch-lint: [no-experimental-imports] package "example/alpha" imports "example/alpha/experimental"`
