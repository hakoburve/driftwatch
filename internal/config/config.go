package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level driftwatch configuration.
type Config struct {
	Version  string    `yaml:"version"`
	Services []Service `yaml:"services"`
}

// Service describes a single service to watch for drift.
type Service struct {
	Name       string            `yaml:"name"`
	Type       string            `yaml:"type"` // e.g. "kubernetes", "docker", "env"
	Source     string            `yaml:"source"`     // path or URL to declared state
	Target     string            `yaml:"target"`     // deployed target identifier
	Labels     map[string]string `yaml:"labels"`
	IgnoreKeys []string          `yaml:"ignore_keys"`
}

// Load reads and parses a driftwatch config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// validate performs basic sanity checks on the loaded configuration.
func (c *Config) validate() error {
	if len(c.Services) == 0 {
		return fmt.Errorf("at least one service must be defined")
	}

	seen := make(map[string]bool)
	for i, svc := range c.Services {
		if svc.Name == "" {
			return fmt.Errorf("service[%d]: name is required", i)
		}
		if seen[svc.Name] {
			return fmt.Errorf("service[%d]: duplicate service name %q", i, svc.Name)
		}
		seen[svc.Name] = true

		if svc.Type == "" {
			return fmt.Errorf("service %q: type is required", svc.Name)
		}
		if svc.Source == "" {
			return fmt.Errorf("service %q: source is required", svc.Name)
		}
		if svc.Target == "" {
			return fmt.Errorf("service %q: target is required", svc.Name)
		}
	}

	return nil
}
