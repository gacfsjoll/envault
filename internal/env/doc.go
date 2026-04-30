// Package env provides utilities for parsing and writing .env files
// used by envault to sync secrets between the local filesystem and
// the encrypted vault storage.
//
// # Parsing
//
// The Parse function reads a .env file and returns a slice of Entry
// values. It handles the following syntax:
//
//   - KEY=value
//   - KEY="quoted value"
//   - KEY='single quoted value'
//   - Inline comments after an unquoted value (e.g. KEY=value # comment)
//   - Full-line comments starting with #
//   - Blank lines (ignored)
//
// # Writing
//
// The Write function serialises a slice of Entry values back to a
// .env file. Values that contain spaces or special characters are
// automatically wrapped in double quotes. The output file is created
// with permissions 0600 to avoid leaking secrets to other users on
// the same machine.
package env
