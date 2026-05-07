// Package cli implements the command-line interface handlers for envault.
//
// Each Run* function corresponds to a top-level subcommand:
//
//   - RunInit   – initialise a new envault project in a directory
//   - RunPush   – encrypt local .env variables into the vault
//   - RunPull   – decrypt vault secrets back into the local .env file
//   - RunList   – display keys (and optionally values) stored in the vault
//   - RunRotate – re-encrypt the vault with a new passphrase
//
// All functions accept a directory path so they can be tested without
// changing the process working directory.
package cli
