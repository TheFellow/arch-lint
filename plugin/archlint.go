package archlint

import (
	"golang.org/x/tools/go/analysis"

	archlint "github.com/TheFellow/arch-lint/pkg/analysis"
)

// New creates a new arch-lint plugin for golangci-lint module plugin system.
func New(settings any) *Plugin {
	return &Plugin{settings: settings}
}

type Plugin struct {
	settings any
}

func (p *Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{archlint.Analyzer}, nil
}

func (p *Plugin) GetLoadMode() string {
	return "syntax"
}
