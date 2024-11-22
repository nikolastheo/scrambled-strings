package tests

import (
	"scrambled-strings/pkg/matcher"
	"testing"
)

func TestCanonicalForm(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"axpaj", "aajpx"},
		{"apxaj", "aajpx"},
		{"dnrbt", "dbnrt"},
		{"pjxdn", "pdjnx"},
		{"a", "a"},
		{"ab", "ab"},
	}

	for _, test := range tests {
		result := matcher.CanonicalForm(test.input)
		if result != test.expected {
			t.Errorf("Expected %s but got %s", test.expected, result)
		}
	}
}
