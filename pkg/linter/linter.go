package linter

import (
	"bufio"
	"fmt"
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

	var violations []Violation
	for _, spec := range cfg.Specs {
		report("spec: %s\n", spec.Name)

		for _, pkg := range pkgs {
			currentPkg := strings.TrimPrefix(pkg.PkgPath, moduleName+"/")

			// Get packages described by the spec include/exclude patterns
			included := false
			for _, includePattern := range spec.Packages.Include {
				if ok, _ := doublestar.Match(includePattern, currentPkg); ok {
					included = true
					break
				}
			}
			if !included {
				continue
			}
			excluded := false
			for _, excludePattern := range spec.Packages.Exclude {
				if ok, _ := doublestar.Match(excludePattern, currentPkg); ok {
					excluded = true
					break
				}
			}
			if excluded {
				continue
			}

			// Validate imports of the current package
			report("  pkg: %q\n", currentPkg)
			for importPath := range pkg.Imports {
				importedPkg := strings.TrimPrefix(importPath, moduleName+"/")
				report("    import: %q\n", importedPkg)
				if v := CheckImport(spec, currentPkg, importedPkg); v != nil {
					violations = append(violations, *v)
				}
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
