package linter

import "github.com/TheFellow/arch-lint/pkg/config"

// CheckImport evaluates whether importedPkg is forbidden for currentPkg
// under the given spec's rules. Returns a *Violation if forbidden, nil otherwise.
func CheckImport(spec config.Spec, currentPkg, importedPkg string) *Violation {
	forbidden := false
	for _, forbid := range spec.Rules.Forbid {
		capturedVars, ok := MatchPattern(forbid, importedPkg)
		if !ok {
			continue
		}
		forbidden = true
		// Exceptions are spec-wide: any matching forbid may supply their captures.
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
	}
	if !forbidden {
		return nil
	}

	return &Violation{
		Rule:    spec.Name,
		Package: currentPkg,
		Import:  importedPkg,
	}
}
