package env

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}
	return path
}

func TestParseSimpleKeyValue(t *testing.T) {
	path := writeTempEnvFile(t, "FOO=bar\nBAZ=qux\n")
	entries, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "FOO" || entries[0].Value != "bar" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
	if entries[1].Key != "BAZ" || entries[1].Value != "qux" {
		t.Errorf("unexpected entry: %+v", entries[1])
	}
}

func TestParseIgnoresComments(t *testing.T) {
	path := writeTempEnvFile(t, "# this is a comment\nKEY=value\n")
	entries, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Key != "KEY" || entries[0].Value != "value" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestParseQuotedValues(t *testing.T) {
	path := writeTempEnvFile(t, `GREETING="hello world"` + "\n")
	entries, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Value != "hello world" {
		t.Errorf("expected 'hello world', got %q", entries[0].Value)
	}
}

func TestParseEmptyFile(t *testing.T) {
	path := writeTempEnvFile(t, "")
	entries, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestWriteRoundTrip(t *testing.T) {
	original := writeTempEnvFile(t, "API_KEY=secret123\nDEBUG=true\n")
	entries, err := Parse(original)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	out := filepath.Join(t.TempDir(), ".env.out")
	if err := Write(out, entries); err != nil {
		t.Fatalf("write error: %v", err)
	}

	parsed, err := Parse(out)
	if err != nil {
		t.Fatalf("re-parse error: %v", err)
	}
	if len(parsed) != len(entries) {
		t.Fatalf("entry count mismatch: want %d got %d", len(entries), len(parsed))
	}
	for i, e := range entries {
		if parsed[i].Key != e.Key || parsed[i].Value != e.Value {
			t.Errorf("entry %d mismatch: want %+v got %+v", i, e, parsed[i])
		}
	}
}

func TestWriteFilePermissions(t *testing.T) {
	entries := []Entry{{Key: "SECRET", Value: "abc"}}
	out := filepath.Join(t.TempDir(), ".env")
	if err := Write(out, entries); err != nil {
		t.Fatalf("write error: %v", err)
	}
	info, err := os.Stat(out)
	if err != nil {
		t.Fatalf("stat error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected permissions 0600, got %v", info.Mode().Perm())
	}
}
