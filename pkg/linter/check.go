package linter

import "github.com/TheFellow/arch-lint/pkg/config"

// CheckImport evaluates whether importedPkg is forbidden for currentPkg
// under the given spec's rules. Returns a *Violation if forbidden, nil otherwise.
func CheckImport(spec config.Spec, currentPkg, importedPkg string) *Violation {
	var capturedVars map[string]string
	forbidden := false
	for _, pat := range spec.Rules.Forbid {
		if vars, ok := MatchPattern(pat, importedPkg); ok {
			capturedVars = vars
			forbidden = true
			break
		}
	}
	if !forbidden {
		return nil
	}

	for _, pat := range spec.Rules.Except {
		if ExceptRegex(pat, currentPkg, capturedVars) {
			return nil
		}
	}

	for _, pat := range spec.Rules.Exempt {
		if ExceptRegex(pat, importedPkg, capturedVars) {
			return nil
		}
	}

	return &Violation{
		Rule:    spec.Name,
		Package: currentPkg,
		Import:  importedPkg,
	}
}
