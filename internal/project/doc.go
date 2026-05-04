// Package project provides per-project configuration for envault.
//
// Each project that uses envault stores a small JSON file (.envault.json)
// in its root directory. This file records:
//
//   - The logical vault name used to look up the encrypted storage file.
//   - The relative path to the .env file managed by envault (defaults to ".env").
//   - An optional default profile name (e.g. "development", "staging").
//
// Typical usage:
//
//	cfg, err := project.Load(".")
//	if errors.Is(err, project.ErrNoConfig) {
//		// prompt user to run `envault init`
//	}
//
// The config file is intentionally human-readable so developers can inspect
// and commit it to version control alongside their project.
package project
