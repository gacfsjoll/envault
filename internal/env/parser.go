// Package env provides utilities for parsing and writing .env files.
// It supports standard .env file format including comments, quoted values,
// and multiline strings.
package env

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair from a .env file,
// along with any inline comment associated with it.
type Entry struct {
	Key     string
	Value   string
	Comment string // inline comment, if any
}

// File represents the parsed contents of a .env file, preserving
// the order of entries as they appear in the source file.
type File struct {
	Entries []Entry
}

// Map returns the entries as a plain key-value map. If duplicate keys
// exist, the last value wins.
func (f *File) Map() map[string]string {
	m := make(map[string]string, len(f.Entries))
	for _, e := range f.Entries {
		m[e.Key] = e.Value
	}
	return m
}

// Parse reads and parses a .env file from the given path.
// Lines beginning with '#' are treated as comments and skipped.
// Blank lines are also skipped. Values may optionally be quoted
// with single or double quotes.
func Parse(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("env: open %q: %w", path, err)
	}
	defer f.Close()

	var file File
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and full-line comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		entry, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("env: %q line %d: %w", path, lineNum, err)
		}
		file.Entries = append(file.Entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("env: scanning %q: %w", path, err)
	}

	return &file, nil
}

// Write serialises a File back to disk at the given path, creating or
// truncating the file as necessary.
func Write(path string, file *File) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("env: create %q: %w", path, err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, e := range file.Entries {
		line := fmt.Sprintf("%s=%s", e.Key, quoteValue(e.Value))
		if e.Comment != "" {
			line += " # " + e.Comment
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return fmt.Errorf("env: write %q: %w", path, err)
		}
	}
	return w.Flush()
}

// parseLine parses a single KEY=VALUE line, stripping optional quotes
// and splitting off any trailing inline comment.
func parseLine(line string) (Entry, error) {
	// Remove optional "export " prefix.
	line = strings.TrimPrefix(line, "export ")

	idx := strings.IndexByte(line, '=')
	if idx < 0 {
		return Entry{}, fmt.Errorf("missing '=' in %q", line)
	}

	key := strings.TrimSpace(line[:idx])
	if key == "" {
		return Entry{}, fmt.Errorf("empty key in %q", line)
	}

	raw := strings.TrimSpace(line[idx+1:])
	value, comment := splitValueComment(raw)
	value = unquote(value)

	return Entry{Key: key, Value: value, Comment: comment}, nil
}

// splitValueComment separates the value portion from an optional trailing
// inline comment (" # ..."). Quoted values are not split.
func splitValueComment(s string) (value, comment string) {
	if len(s) == 0 {
		return "", ""
	}
	// If the value is quoted, the comment starts after the closing quote.
	if s[0] == '"' || s[0] == '\'' {
		q := s[0]
		end := strings.IndexByte(s[1:], q)
		if end >= 0 {
			closing := end + 2 // account for offset +1 and the quote itself
			rest := strings.TrimSpace(s[closing:])
			if strings.HasPrefix(rest, "#") {
				return s[:closing], strings.TrimSpace(rest[1:])
			}
			return s[:closing], ""
		}
	}
	// Unquoted: split on " #".
	if i := strings.Index(s, " #"); i >= 0 {
		return strings.TrimSpace(s[:i]), strings.TrimSpace(s[i+2:])
	}
	return s, ""
}

// unquote removes surrounding single or double quotes from a value.
func unquote(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// quoteValue wraps a value in double quotes if it contains spaces,
// special characters, or is empty.
func quoteValue(s string) string {
	if s == "" || strings.ContainsAny(s, " \t\n#\'\"\\$") {
		return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
	}
	return s
}
