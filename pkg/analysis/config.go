package analysis

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/TheFellow/arch-lint/pkg/config"
)

var configCache sync.Map

type cachedConfig struct {
	cfg *config.Config
	err error
}

func loadConfigCached(startDir, flagPath string) (*config.Config, error) {
	if flagPath != "" {
		absPath, _ := filepath.Abs(flagPath)
		if cached, ok := configCache.Load(absPath); ok {
			c := cached.(cachedConfig)
			return c.cfg, c.err
		}

		cfg, err := config.Load(absPath)
		configCache.Store(absPath, cachedConfig{cfg, err})
		return cfg, err
	}

	configPath, err := resolveConfigPath(startDir)
	if err != nil {
		return nil, err
	}
	if cached, ok := configCache.Load(configPath); ok {
		c := cached.(cachedConfig)
		return c.cfg, c.err
	}

	cfg, loadErr := config.Load(configPath)
	configCache.Store(configPath, cachedConfig{cfg, loadErr})
	return cfg, loadErr
}

func resolveConfigPath(startDir string) (string, error) {
	current := startDir
	for {
		candidate := filepath.Join(current, ".arch-lint.yml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", fmt.Errorf("arch-lint: no .arch-lint.yml found (use -config flag)")
		}
		current = parent
	}
}
