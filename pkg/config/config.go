package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Type       string   `yaml:"type"`
	Name       string   `yaml:"name"`
	Packages   []string `yaml:"packages"`   // go list patterns to target specific packages
	Forbidden  []string `yaml:"forbidden"`  // forbidden import package globs
	Exceptions []string `yaml:"exceptions"` // exceptions to forbidden imports
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	for _, r := range cfg.Rules {
		if len(r.Packages) == 0 {
			return nil, fmt.Errorf("rule '%s' must specify 'packages'", r.Name)
		}
	}
	return &cfg, nil
}
