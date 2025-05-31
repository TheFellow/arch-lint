package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	IncludeTests bool   `yaml:"include_tests"`
	Specs        []Spec `yaml:"specs"`
}

type Spec struct {
	Name     string   `yaml:"name"`
	Packages Packages `yaml:"packages"`
	Rules    Rules    `yaml:"rules"`
}

type Rules struct {
	Forbid []string `yaml:"forbid"`
	Except []string `yaml:"except"`
	Exempt []string `yaml:"exempt"`
}

type Files struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}

type Packages struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
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
	if len(cfg.Specs) == 0 {
		return nil, fmt.Errorf("config must contain at least one spec")
	}
	for _, r := range cfg.Specs {
		if len(r.Packages.Include) == 0 {
			return nil, fmt.Errorf("rule '%s' must specify 'packages'", r.Name)
		}
		if len(r.Rules.Forbid) == 0 {
			return nil, fmt.Errorf("rule '%s' must specify 'forbid' rules", r.Name)
		}
	}
	return &cfg, nil
}
