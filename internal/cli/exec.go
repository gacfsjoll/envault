package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/envault/envault/internal/env"
	"github.com/envault/envault/internal/project"
	"github.com/envault/envault/internal/storage"
)

// RunExec loads secrets from the vault and injects them as environment variables
// into the provided command, then executes it.
func RunExec(passphrase string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command provided")
	}

	cfg, err := project.Load(".")
	if err != nil {
		return fmt.Errorf("failed to load project config: %w", err)
	}

	vault, err := storage.New(cfg.VaultName, passphrase)
	if err != nil {
		return fmt.Errorf("failed to open vault: %w", err)
	}

	secrets := vault.All()

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start with a copy of the current environment.
	cmd.Env = os.Environ()

	// Overlay vault secrets, converting to KEY=VALUE pairs.
	for k, v := range secrets {
		cmd.Env = env.SetEnvVar(cmd.Env, k, v)
	}

	if err := cmd.Run(); err != nil {
		// Propagate the exit code naturally.
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}
