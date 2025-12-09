package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"gopkg.in/yaml.v3"
)

// User represents a GitHub user to monitor
type User struct {
	Name  string `yaml:"name" toml:"name" hcl:"name"`
	Token string `yaml:"token" toml:"token" hcl:"token"`
}

// Config represents the application configuration
type Config struct {
	Users      []User `yaml:"users" toml:"users" hcl:"user,block"`
	ListenAddr string `yaml:"listen_addr,omitempty" toml:"listen_addr,omitempty" hcl:"listen_addr,optional"`
	MetricsPath string `yaml:"metrics_path,omitempty" toml:"metrics_path,omitempty" hcl:"metrics_path,optional"`
	PollInterval int   `yaml:"poll_interval,omitempty" toml:"poll_interval,omitempty" hcl:"poll_interval,optional"`
}

// LoadConfig loads configuration from a file (YAML, TOML, or HCL based on extension)
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(path))

	var cfg Config

	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &cfg)
	case ".toml":
		err = toml.Unmarshal(data, &cfg)
	case ".hcl":
		err = hclsimple.Decode(path, data, nil, &cfg)
	default:
		return nil, fmt.Errorf("unsupported config file format: %s (supported: .yaml, .yml, .toml, .hcl)", ext)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = ":9101"
	}
	if cfg.MetricsPath == "" {
		cfg.MetricsPath = "/metrics"
	}
	if cfg.PollInterval == 0 {
		cfg.PollInterval = 60 // Default to 60 seconds
	}

	// Validate
	if len(cfg.Users) == 0 {
		return nil, fmt.Errorf("no users defined in config")
	}

	for i, user := range cfg.Users {
		if user.Name == "" {
			return nil, fmt.Errorf("user at index %d has no name", i)
		}
		if user.Token == "" {
			return nil, fmt.Errorf("user %s has no token", user.Name)
		}
	}

	return &cfg, nil
}
