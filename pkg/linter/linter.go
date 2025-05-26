package linter

import (
	"bufio"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/TheFellow/arch-lint/pkg/config"
	"github.com/bmatcuk/doublestar/v4"
)

var Processed strings.Builder

func report(str string, args ...any) {
	Processed.WriteString(fmt.Sprintf(str, args...))
}

// Run enforces forbidden import rules by analyzing files specified by glob patterns
func Run(cfg *config.Config) ([]Violation, error) {

	moduleName, err := getModuleName()
	if err != nil {
		return nil, err
	}

	pkgs, err := loadPackages(cfg)
	if err != nil {
		return nil, err
	}

	fileToPackagePath, err := mapFilesToPackages(moduleName, pkgs)
	if err != nil {
		return nil, err
	}

	var violations []Violation
	for _, spec := range cfg.Specs {
		report("spec: %s\n", spec.Name)

		files, err := getFilesToCheck(fileToPackagePath, spec)
		if err != nil {
			return nil, err
		}
		report("  %d files\n", len(files))

		for _, file := range files {
			report("  file: %q\n", file)

			// Parse the file to extract imports and package name
			// TODO: Pull this from analysis packages instead of reparsing the file
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, file, nil, parser.ImportsOnly)
			if err != nil {
				return nil, fmt.Errorf("failed to parse file %s: %w", file, err)
			}

			// Get the full package path from the map
			packagePath, ok := fileToPackagePath[file]
			if !ok {
				return nil, fmt.Errorf("spec: %s: package path not found for file: %s", spec.Name, file)
			}
			report("   pkg: %q\n", packagePath)

			for _, imp := range node.Imports {
				// Extract the relative import path
				importPath := strings.Trim(imp.Path.Value, `"`)
				importPath = strings.TrimPrefix(importPath, moduleName+"/")
				report("    import: %q\n", importPath)

				// Check imports for forbidden rules
				var capturedVars map[string]string
				forbidden := false
				for _, pat := range spec.Rules.Forbid {
					if vars, ok := matchPattern(pat, importPath); ok {
						report("      forbid: %q\n", pat)
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
				for _, pat := range spec.Rules.Except {
					if exceptRegex(pat, packagePath, capturedVars) {
						report("      except: %q\n", pat)
						exception = true
						break
					}
				}
				if exception {
					continue
				}

				// If the import is forbidden and not in exceptions, add a violation
				violations = append(violations, Violation{
					Rule:    spec.Name,
					File:    file,
					Package: packagePath,
					Import:  importPath,
				})
			}
		}

	}

	return violations, nil
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

// loadPackages pulls in all package information using the analysis package
func loadPackages(cfg *config.Config) ([]*packages.Package, error) {
	// Load package information using the analysis package
	cfgs := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedImports,
	}
	if cfg.IncludeTests {
		cfgs.Tests = true
		cfgs.Mode |= packages.NeedForTest
	}
	pkgs, err := packages.Load(cfgs, "./...")
	if err != nil {
		return nil, fmt.Errorf("failed to load packages: %w", err)
	}
	return pkgs, nil
}

// mapFilesToPackages maps file paths to their relative package names. So
//
//	"github.com/TheFellow/example/alpha/tester.go" -> "example/alpha"
func mapFilesToPackages(moduleName string, pkgs []*packages.Package) (map[string]string, error) {
	// Map file paths to their full package paths
	fileToPackagePath := make(map[string]string)
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	for _, pkg := range pkgs {
		for _, fileAbs := range pkg.GoFiles {
			file := strings.TrimPrefix(fileAbs, wd+"/")
			pkgPath := strings.TrimPrefix(pkg.PkgPath, moduleName+"/")
			fileToPackagePath[file] = pkgPath
		}
	}
	return fileToPackagePath, nil
}

// getFilesToCheck returns a list of files specified by Include not excluded by Exclude globs
func getFilesToCheck(fileToPackagePath map[string]string, spec config.Spec) ([]string, error) {
	var filesToCheck []string

	var includedFiles []string
	for _, includePattern := range spec.Packages.Include {
		for file, pkg := range fileToPackagePath {
			ok, err := doublestar.Match(includePattern, pkg)
			if err != nil {
				return nil, fmt.Errorf("failed to match include pattern %s: %w", includePattern, err)
			}
			if ok {
				includedFiles = append(includedFiles, file)
			}
		}
	}

	var excludedFiles []string
	for _, excludePattern := range spec.Packages.Exclude {
		for file, pkg := range fileToPackagePath {
			ok, err := doublestar.Match(excludePattern, pkg)
			if err != nil {
				return nil, fmt.Errorf("failed to match include pattern %s: %w", excludePattern, err)
			}
			if ok {
				excludedFiles = append(excludedFiles, file)
			}
		}
	}

	excludedSet := make(map[string]struct{})
	for _, file := range excludedFiles {
		excludedSet[file] = struct{}{}
	}

	for _, file := range includedFiles {
		if _, excluded := excludedSet[file]; !excluded {
			filesToCheck = append(filesToCheck, file)
		}
	}

	return filesToCheck, nil
}
