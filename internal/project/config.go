// Package project manages per-project envault configuration,
// linking a project directory to a named vault and environment profiles.
package project

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const configFileName = ".envault.json"

// ErrNoConfig is returned when no project config file is found in the directory.
var ErrNoConfig = errors.New("no envault config found in directory")

// Config holds the project-level envault settings stored in .envault.json.
type Config struct {
	// VaultName is the logical name of the vault associated with this project.
	VaultName string `json:"vault_name"`
	// DefaultProfile is the env profile loaded by default (e.g. "development").
	DefaultProfile string `json:"default_profile,omitempty"`
	// EnvFile is the relative path to the .env file managed by envault.
	EnvFile string `json:"env_file"`
}

// Load reads the .envault.json config from the given directory.
// Returns ErrNoConfig if the file does not exist.
func Load(dir string) (*Config, error) {
	path := filepath.Join(dir, configFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ErrNoConfig
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.VaultName == "" {
		return nil, errors.New("config missing required field: vault_name")
	}
	if cfg.EnvFile == "" {
		cfg.EnvFile = ".env"
	}
	return &cfg, nil
}

// Save writes the Config as .envault.json into the given directory.
// The file is created with mode 0644.
func Save(dir string, cfg *Config) error {
	if cfg.VaultName == "" {
		return errors.New("vault_name must not be empty")
	}
	if cfg.EnvFile == "" {
		cfg.EnvFile = ".env"
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(dir, configFileName)
	return os.WriteFile(path, append(data, '\n'), 0644)
}
