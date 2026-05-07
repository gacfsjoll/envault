package env

import (
	"testing"
)

func TestSetEnvVarAppendsNew(t *testing.T) {
	env := []string{"FOO=bar", "BAZ=qux"}
	result := SetEnvVar(env, "NEW_KEY", "new_value")
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result[2] != "NEW_KEY=new_value" {
		t.Errorf("unexpected last entry: %s", result[2])
	}
}

func TestSetEnvVarReplacesExisting(t *testing.T) {
	env := []string{"FOO=bar", "BAZ=qux"}
	result := SetEnvVar(env, "FOO", "replaced")
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0] != "FOO=replaced" {
		t.Errorf("expected FOO=replaced, got %s", result[0])
	}
}

func TestSetEnvVarEmptySlice(t *testing.T) {
	result := SetEnvVar(nil, "KEY", "val")
	if len(result) != 1 || result[0] != "KEY=val" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestToEnvSlice(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2"}
	slice := ToEnvSlice(vars)
	if len(slice) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(slice))
	}
	seen := make(map[string]bool)
	for _, e := range slice {
		seen[e] = true
	}
	if !seen["A=1"] || !seen["B=2"] {
		t.Errorf("unexpected entries: %v", slice)
	}
}

func TestToEnvSliceEmpty(t *testing.T) {
	slice := ToEnvSlice(map[string]string{})
	if len(slice) != 0 {
		t.Errorf("expected empty slice, got %v", slice)
	}
}
