package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/envault/envault/internal/project"
)

const defaultEnvFile = ".env"
const defaultVaultFile = ".envault"

// InitOptions holds the configuration for the init command.
type InitOptions struct {
	VaultName string
	EnvFile   string
	Force     bool
}

// RunInit initialises a new envault project in the current directory.
// It creates a .envault.json config file and refuses to overwrite an
// existing one unless Force is set.
func RunInit(opts InitOptions) error {
	if opts.VaultName == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("could not determine working directory: %w", err)
		}
		opts.VaultName = filepath.Base(cwd)
	}

	if opts.EnvFile == "" {
		opts.EnvFile = defaultEnvFile
	}

	cfgPath := project.ConfigFileName
	if !opts.Force {
		if _, err := os.Stat(cfgPath); err == nil {
			return errors.New("project already initialised; use --force to overwrite")
		}
	}

	cfg := &project.Config{
		VaultName: opts.VaultName,
		EnvFile:   opts.EnvFile,
	}

	if err := project.Save(cfg); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("Initialised envault project\n")
	fmt.Printf("  vault name : %s\n", cfg.VaultName)
	fmt.Printf("  env file   : %s\n", cfg.EnvFile)
	fmt.Printf("  config     : %s\n", cfgPath)
	return nil
}
