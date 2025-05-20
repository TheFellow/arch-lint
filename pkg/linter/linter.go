package linter

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/TheFellow/go-arch-lint/pkg/config"
	"github.com/bmatcuk/doublestar/v4"
)

var debug = false

func log(str string, args ...any) {
	if debug {
		fmt.Printf(str, args...)
	}
}

// Run enforces forbidden import rules by analyzing files specified by glob patterns
func Run(cfg *config.Config) ([]Violation, error) {
	var violations []Violation

	for _, rule := range cfg.Rules {
		log("Checking rule: %s\n", rule.Name)

		for _, pkgPattern := range rule.Packages {
			log("--Checking package pattern: %s\n", pkgPattern)

			// Resolve file paths using the glob pattern
			files, err := doublestar.Glob(os.DirFS("."), pkgPattern)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve glob pattern %s: %w", pkgPattern, err)
			}
			//fmt.Println(os.Getwd())
			log("  --Found %d files\n", len(files))

			for _, file := range files {
				log("  --Checking file: %s\n", file)

				// Parse the file to extract imports and package name
				fset := token.NewFileSet()
				node, err := parser.ParseFile(fset, file, nil, parser.ImportsOnly)
				if err != nil {
					return nil, fmt.Errorf("failed to parse file %s: %w", file, err)
				}

				// Construct the package path from the file path
				packagePath := strings.TrimPrefix(filepath.Dir(file), "./")
				packagePath = strings.ReplaceAll(packagePath, string(os.PathSeparator), "/")

				for _, imp := range node.Imports {
					importPath := strings.Trim(imp.Path.Value, `"`)
					log("    --Found import: %s\n", importPath)

					// Check imports for forbidden rules
					forbidden := false
					for _, pat := range rule.Forbidden {
						if ok, _ := doublestar.Match(pat, importPath); ok {
							log("    --Forbidden by: %s\n", pat)
							forbidden = true
							break
						}
					}
					if !forbidden {
						continue
					}

					// Check if the source package is in exceptions
					exception := false
					for _, pat := range rule.Exceptions {
						if ok, _ := doublestar.Match(pat, packagePath); ok {
							log("    --Exempted by: %s\n", pat)
							exception = true
							break
						}
					}
					if exception {
						continue
					}

					// If the import is forbidden and not in exceptions, add a violation
					violations = append(violations, Violation{
						File:   file,
						Import: importPath,
						Rule:   rule.Name,
					})
				}
			}
		}
	}

	return violations, nil
}
