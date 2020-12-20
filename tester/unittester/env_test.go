package unittester

import (
	"testing"
)

func TestReplaceWithEnvValue(t *testing.T) {
	tests := []struct {
		input    string
		env      map[string]string
		expected string
	}{
		{"empty", map[string]string{}, "empty"},
		{"em#p#ty", map[string]string{}, "em#p#ty"},
		{"em#p#ty", map[string]string{"a": "b"}, "em#p#ty"},
		{"em#p#ty", map[string]string{"p": "p"}, "empty"},
		{"em#p#ty#a#", map[string]string{"p": "p"}, "empty#a#"},
		{"em#p#ty#a#", map[string]string{"p": "p", "a": "p"}, "emptyp"},
	}

	for i, tt := range tests {
		result := ReplaceStringWithEnvValue(tt.input, tt.env)
		if result != tt.expected {
			t.Fatalf("%d : failed got %s exp %s", i, result, tt.expected)
		}
		asBytes := []byte(tt.input)
		bytesResult := ReplaceWithEnvValue(asBytes, tt.env)
		if string(bytesResult) != tt.expected {
			t.Fatalf("%d : failed got %s exp %s", i, string(bytesResult), tt.expected)
		}
	}
}
