package analysis

import (
	"path/filepath"
	"sync"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	configCache = sync.Map{}
	configFlag = ""
	t.Cleanup(func() {
		configCache = sync.Map{}
		configFlag = ""
	})

	if err := Analyzer.Flags.Set("config",
		filepath.Join(testdata, "src", "example", ".arch-lint.yml")); err != nil {
		t.Fatalf("set config flag: %v", err)
	}

	analysistest.Run(t, testdata, Analyzer, "example/...")
}
