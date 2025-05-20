package linter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/TheFellow/go-arch-lint/pkg/config"
	"github.com/bmatcuk/doublestar/v4"
)

// GoPackage mirrors `go list -json` output
type GoPackage struct {
	ImportPath string   `json:"ImportPath"`
	Imports    []string `json:"Imports"`
}

// Violation represents a rule violation
type Violation struct {
	Pkg     string
	Import  string
	Rule    string
	Message string
}

func (v Violation) String() string {
	return fmt.Sprintf("%s imports %s [%s]: %s", v.Pkg, v.Import, v.Rule, v.Message)
}

// hasCommonAncestor checks if two import paths share a directory prefix
func hasCommonAncestor(a, b string) bool {
	pa := strings.Split(a, "/")
	pb := strings.Split(b, "/")
	min := len(pa)
	if len(pb) < min {
		min = len(pb)
	}
	for i := 0; i < min; i++ {
		if pa[i] != pb[i] {
			return i > 0
		}
	}
	return true
}

// Run invokes `go list -json` on user-specified package patterns and enforces forbidden import rules
func Run(cfg *config.Config) ([]Violation, error) {
	// compile unique package patterns
	patterns := []string{""}
	for _, rule := range cfg.Rules {
		patterns = append(patterns, rule.Packages...)
	}
	// exec go list
	args := append([]string{"list", "-json"}, patterns...)
	cmd := exec.Command("go", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run go list: %w", err)
	}

	// decode packages
	dec := json.NewDecoder(bytes.NewReader(out))
	pkgMap := make(map[string]*GoPackage)
	for {
		var gp GoPackage
		if err := dec.Decode(&gp); err != nil {
			break
		}
		pkgMap[gp.ImportPath] = &gp
	}

	var violations []Violation
	for _, rule := range cfg.Rules {
		for _, pkgPattern := range rule.Packages {
			for importPath, gp := range pkgMap {
				if ok, _ := doublestar.Match(pkgPattern, importPath); !ok {
					continue
				}

				for _, imp := range gp.Imports {
					// forbidden?
					forbidden := false
					for _, pat := range rule.Forbidden {
						if ok, _ := doublestar.Match(pat, imp); ok {
							forbidden = true
							break
						}
					}
					if !forbidden {
						continue
					}

					// exceptions
					exc := false
					for _, pat := range rule.Exceptions {
						if ok, _ := doublestar.Match(pat, imp); ok {
							exc = true
							break
						}
					}
					if exc {
						continue
					}

					// allow same ancestor
					if hasCommonAncestor(importPath, imp) {
						continue
					}

					violations = append(violations, Violation{Pkg: importPath, Import: imp, Rule: rule.Name, Message: "forbidden import"})
				}
			}
		}
	}
	return violations, nil
}
