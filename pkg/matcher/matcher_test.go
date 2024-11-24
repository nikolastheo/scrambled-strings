package matcher

import (
	"testing"
)

// TestCanonicalForm validates the canonicalForm function.
func TestCanonicalForm(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"axpaj", "aapxj"},
		{"apxaj", "aapxj"},
		{"dnrbt", "dbnrt"},
		{"pjxdn", "pdjxn"},
		{"a", "a"},
		{"ab", "ab"},
	}

	for _, test := range tests {
		result := CanonicalForm(test.input)
		if result != test.expected {
			t.Errorf("Expected %s but got %s", test.expected, result)
		}
	}
}

// TestPrecomputeCanonicalForms validates dictionary preprocessing.
func TestPrecomputeCanonicalForms(t *testing.T) {
	dictionary := []string{"axpaj", "apxaj", "dnrbt", "pjxdn", "abd"}
	expected := map[string]string{
		"axpaj": "aapxj",
		"apxaj": "aapxj",
		"dnrbt": "dbnrt",
		"pjxdn": "pdjxn",
		"abd":   "abd",
	}

	result := PrecomputeCanonicalForms(dictionary)

	if len(result) != len(expected) {
		t.Errorf("Expected %d entries but got %d", len(expected), len(result))
	}

	for word, canonical := range expected {
		if result[word] != canonical {
			t.Errorf("Expected %s -> %s but got %s", word, canonical, result[word])
		}
	}
}

// TestCountMatches validates the matching logic.
func TestCountMatches(t *testing.T) {
	dictionary := []string{"axpaj", "apxaj", "dnrbt", "pjxdn", "abd"}
	precomputed := PrecomputeCanonicalForms(dictionary)
	input := "aapxjdnrbtvldptfzbbdbbzxtndrvjblnzjfpvhdhhpxjdnrbt"

	expected := 4
	result := CountMatches(precomputed, input)

	if result != expected {
		t.Errorf("Expected %d matches but got %d", expected, result)
	}
}
