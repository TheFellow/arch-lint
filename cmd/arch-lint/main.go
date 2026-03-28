package main

import (
	archlint "github.com/TheFellow/arch-lint/pkg/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(archlint.Analyzer)
}
