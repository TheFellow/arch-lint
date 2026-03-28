package analysis

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/TheFellow/arch-lint/pkg/config"
	"github.com/TheFellow/arch-lint/pkg/linter"
	"github.com/bmatcuk/doublestar/v4"
)

var Analyzer = &analysis.Analyzer{
	Name: "archlint",
	Doc:  "enforces architectural import boundaries via configurable rules",
	Run:  run,
}

func init() {
	Analyzer.Flags.StringVar(&configFlag, "config", "",
		"path to arch-lint config file (default: walk up for .arch-lint.yml)")
}

var configFlag string

func run(pass *analysis.Pass) (interface{}, error) {
	dir := projectDir(pass)
	cfg, err := loadConfigCached(dir, configFlag)
	if err != nil {
		return nil, err
	}

	modulePath := resolveModulePath(pass)
	currentPkg := strings.TrimPrefix(pass.Pkg.Path(), modulePath+"/")

	if !cfg.IncludeTests && isTestPackage(pass) {
		return nil, nil
	}

	for _, spec := range cfg.Specs {
		if !matchesInclude(spec, currentPkg) || matchesExclude(spec, currentPkg) {
			continue
		}

		for _, file := range pass.Files {
			for _, imp := range file.Imports {
				importPath := strings.Trim(imp.Path.Value, `"`)
				importedPkg := strings.TrimPrefix(importPath, modulePath+"/")

				if v := linter.CheckImport(spec, currentPkg, importedPkg); v != nil {
					pass.Reportf(imp.Pos(), "[%s] forbidden import of %q", spec.Name, importedPkg)
				}
			}
		}
	}

	return nil, nil
}

func matchesInclude(spec config.Spec, pkg string) bool {
	for _, pattern := range spec.Packages.Include {
		if ok, _ := doublestar.Match(pattern, pkg); ok {
			return true
		}
	}
	return false
}

func matchesExclude(spec config.Spec, pkg string) bool {
	for _, pattern := range spec.Packages.Exclude {
		if ok, _ := doublestar.Match(pattern, pkg); ok {
			return true
		}
	}
	return false
}

func isTestPackage(pass *analysis.Pass) bool {
	return strings.HasSuffix(pass.Pkg.Path(), "_test") ||
		strings.HasSuffix(pass.Pkg.Name(), "_test")
}

func resolveModulePath(pass *analysis.Pass) string {
	if pass.Module != nil && pass.Module.Path != "" {
		return pass.Module.Path
	}

	dir := projectDir(pass)
	if mod, err := findGoMod(dir); err == nil {
		return mod
	}
	return ""
}

func projectDir(pass *analysis.Pass) string {
	if len(pass.Files) == 0 {
		return "."
	}

	file := pass.Fset.File(pass.Files[0].Pos())
	if file == nil {
		return "."
	}
	return filepath.Dir(file.Name())
}

func findGoMod(dir string) (string, error) {
	current := dir
	for {
		goModPath := filepath.Join(current, "go.mod")
		if data, err := os.ReadFile(goModPath); err == nil {
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "module ") {
					return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
				}
			}
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", fmt.Errorf("go.mod not found")
		}
		current = parent
	}
}
