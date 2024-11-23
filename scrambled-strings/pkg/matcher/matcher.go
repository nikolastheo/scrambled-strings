package matcher

import (
	"sort"
)

// CanonicalForm computes the canonical version of a word by sorting the middle characters
// and preserving the first and last characters in their original positions.
func CanonicalForm(word string) string {
	// If the word has 2 or fewer characters, return it as-is
	if len(word) <= 2 {
		return word
	}

	// Extract middle characters
	middle := []rune(word[1 : len(word)-1])

	// Sort the middle characters alphabetically
	sort.Slice(middle, func(i, j int) bool {
		return middle[i] < middle[j]
	})

	// Reconstruct the canonical word
	return string(word[0]) + string(middle) + string(word[len(word)-1])
}


// PrecomputeCanonicalForms generates a map of canonical forms for a list of dictionary words.
func PrecomputeCanonicalForms(words []string) map[string]string {
	canonicalMap := make(map[string]string)
	for _, word := range words {
		if _, exists := canonicalMap[word]; !exists {
			canonicalMap[word] = CanonicalForm(word)
		}
	}
	return canonicalMap
}

// CountMatches finds the number of dictionary words (original or scrambled) that match substrings in the input.
func CountMatches(dictionary map[string]string, input string) int {
	count := 0

	for word, canonical := range dictionary {
		wordLen := len(word)

		for i := 0; i <= len(input)-wordLen; i++ {
			substring := input[i : i+wordLen]
			substringCanonical := CanonicalForm(substring)

			if substringCanonical == canonical {
				count++
				break // Stop after the first match for this word
			}
		}
	}

	return count
}