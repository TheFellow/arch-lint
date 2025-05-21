package linter

import (
	"bufio"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/TheFellow/go-arch-lint/pkg/config"
	"github.com/bmatcuk/doublestar/v4"
)

var debug = true

func log(str string, args ...any) {
	if debug {
		fmt.Printf(str, args...)
	}
}

// getModuleName extracts the module name from go.mod
func getModuleName() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", fmt.Errorf("failed to open go.mod: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	return "", fmt.Errorf("module name not found in go.mod")
}

// Run enforces forbidden import rules by analyzing files specified by glob patterns
func Run(cfg *config.Config) ([]Violation, error) {
	var violations []Violation

	// Get the module name from go.mod
	moduleName, err := getModuleName()
	if err != nil {
		return nil, err
	}

	for _, spec := range cfg.Specs {
		log("spec: %s\n", spec.Name)

		files, err := getFilesToCheck(spec.Files)
		if err != nil {
			return nil, err
		}
		log("  %d files\n", len(files))

		for _, file := range files {
			log("  file: %q\n", file)

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
				importPath = strings.TrimPrefix(importPath, moduleName+"/")
				log("    import: %q\n", importPath)

				// Check imports for forbidden rules
				var capturedVars map[string]string
				forbidden := false
				for _, pat := range spec.Rules {
					log("      check: %q\n", pat.Forbid)
					if vars, ok := matchPattern(pat.Forbid, importPath); ok {
						capturedVars = vars
						forbidden = true
						break
					}
				}
				if !forbidden {
					continue
				}

				// Check if the source package is in exceptions
				exception := false
				for _, pat := range spec.Rules {
					exceptPattern := replaceVariables(pat.Except, capturedVars)
					if _, ok := matchPattern(exceptPattern, packagePath); ok {
						log("      exempt: %q\n", pat.Except)
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
					Rule:   spec.Name,
				})
			}
		}

	}

	return violations, nil
}

func getFilesToCheck(files config.Files) ([]string, error) {
	// Resolve include globs
	var includedFiles []string
	for _, includePattern := range files.Include {
		files, err := doublestar.Glob(os.DirFS("."), includePattern)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve include glob pattern %s: %w", includePattern, err)
		}
		includedFiles = append(includedFiles, files...)
	}

	// Resolve exclude globs
	var excludedFiles []string
	for _, excludePattern := range files.Exclude {
		files, err := doublestar.Glob(os.DirFS("."), excludePattern)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve exclude glob pattern %s: %w", excludePattern, err)
		}
		excludedFiles = append(excludedFiles, files...)
	}

	// Filter out excluded files
	excludedSet := make(map[string]struct{})
	for _, file := range excludedFiles {
		excludedSet[file] = struct{}{}
	}

	var filesToCheck []string
	for _, file := range includedFiles {
		if _, excluded := excludedSet[file]; !excluded {
			filesToCheck = append(filesToCheck, file)
		} else {
			log("skip file: %s\n", file)
		}
	}
	return filesToCheck, nil
}
