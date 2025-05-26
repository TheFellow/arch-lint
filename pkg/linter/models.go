package linter

import "fmt"

// GoPackage mirrors `go list -json` output
type GoPackage struct {
	ImportPath string   `json:"ImportPath"`
	Imports    []string `json:"Imports"`
}

// Violation represents a rule violation
type Violation struct {
	File    string
	Package string
	Import  string
	Rule    string
}

func (v Violation) String() string {
	return fmt.Sprintf("arch-lint: [%s] file %q (package %s) imports %q", v.Rule, v.File, v.Package, v.Import)
}
