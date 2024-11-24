package matcher

import (
	"sort"

	"github.com/sirupsen/logrus"
)

// Set up logger
var log = logrus.New()

// CanonicalForm computes the canonical version of a word by sorting the middle characters
// and preserving the first and last characters in their original positions.
func CanonicalForm(word string) string {
	log.WithField("word", word).Debug("Starting computation of canonical form")

	// If the word has 2 or fewer characters, return it as-is
	if len(word) <= 2 {
		log.WithField("canonicalForm", word).Debug("Word is 2 or fewer characters; returning as-is")
		return word
	}

	// Extract middle characters
	middle := []rune(word[1 : len(word)-1])
	log.WithField("middleCharacters", string(middle)).Trace("Extracted middle characters")

	// Sort the middle characters alphabetically
	sort.Slice(middle, func(i, j int) bool {
		return middle[i] < middle[j]
	})
	log.WithField("sortedMiddle", string(middle)).Trace("Sorted middle characters")

	// Reconstruct the canonical word
	canonical := string(word[0]) + string(middle) + string(word[len(word)-1])
	log.WithFields(logrus.Fields{
		"originalWord":  word,
		"canonicalForm": canonical,
	}).Debug("Computed canonical form")
	return canonical
}

// PrecomputeCanonicalForms generates a map of canonical forms for a list of dictionary words.
func PrecomputeCanonicalForms(words []string) map[string]string {
	log.WithField("wordCount", len(words)).Info("Precomputing canonical forms for dictionary words")
	canonicalMap := make(map[string]string)

	for _, word := range words {
		if _, exists := canonicalMap[word]; !exists {
			canonical := CanonicalForm(word)
			canonicalMap[word] = canonical
			log.WithFields(logrus.Fields{
				"word":          word,
				"canonicalForm": canonical,
			}).Trace("Added canonical form to dictionary map")
		} else {
			log.WithField("word", word).Trace("Skipped duplicate word in dictionary")
		}
	}

	log.WithField("canonicalMapSize", len(canonicalMap)).Info("Finished precomputing canonical forms")
	return canonicalMap
}

// CountMatches finds the number of dictionary words (original or scrambled) that match substrings in the input.
func CountMatches(dictionary map[string]string, input string) int {
	log.WithField("inputLength", len(input)).Info("Starting match counting")
	count := 0

	for word, canonical := range dictionary {
		wordLen := len(word)
		log.WithField("word", word).Trace("Processing dictionary word")

		for i := 0; i <= len(input)-wordLen; i++ {
			substring := input[i : i+wordLen]
			substringCanonical := CanonicalForm(substring)
			log.WithFields(logrus.Fields{
				"substring":          substring,
				"substringCanonical": substringCanonical,
			}).Trace("Generated canonical form for substring")

			if substringCanonical == canonical {
				log.WithFields(logrus.Fields{
					"word":      word,
					"canonical": canonical,
					"substring": substring,
				}).Info("Match found")
				count++
				break // Stop after the first match for this word
			}
		}
	}

	log.WithField("totalMatches", count).Info("Finished counting matches")
	return count
}
