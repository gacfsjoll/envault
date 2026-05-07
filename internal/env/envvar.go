package env

import "strings"

// SetEnvVar sets or replaces a KEY=VALUE entry in an environment slice.
// If a variable with the same key already exists it is replaced in-place;
// otherwise the new entry is appended.
func SetEnvVar(environ []string, key, value string) []string {
	prefix := key + "="
	for i, entry := range environ {
		if strings.HasPrefix(entry, prefix) {
			environ[i] = prefix + value
			return environ
		}
	}
	return append(environ, prefix+value)
}

// ToEnvSlice converts a map of key/value pairs to a slice of KEY=VALUE strings
// suitable for use as os/exec Env.
func ToEnvSlice(vars map[string]string) []string {
	result := make([]string, 0, len(vars))
	for k, v := range vars {
		result = append(result, k+"="+v)
	}
	return result
}
