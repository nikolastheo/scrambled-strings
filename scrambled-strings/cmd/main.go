package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"scrambled-strings/pkg/matcher"
)

// readLines reads a file and returns its lines as a slice of strings.
func readLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
func run(dictionaryPath, inputPath string) error {
	// Read and preprocess the dictionary
	dictionaryWords, err := readLines(dictionaryPath)
	if err != nil {
		return fmt.Errorf("error reading dictionary file: %v", err)
	}
	canonicalDictionary := matcher.PrecomputeCanonicalForms(dictionaryWords)

	// Read input strings
	inputStrings, err := readLines(inputPath)
	if err != nil {
		return fmt.Errorf("error reading input file: %v", err)
	}

	// Process each input string and print results
	for i, input := range inputStrings {
		matchCount := matcher.CountMatches(canonicalDictionary, input)
		fmt.Printf("Case #%d: %d\n", i+1, matchCount)
	}

	return nil
}

func main() {
	// Define command-line flags for dictionary and input file paths
	dictionaryPath := flag.String("dictionary", "", "Path to the dictionary file")
	inputPath := flag.String("input", "", "Path to the input file")
	flag.Parse()

	// Ensure both flags are provided
	if *dictionaryPath == "" || *inputPath == "" {
		fmt.Println("Usage: ./scrambled-strings --dictionary <path> --input <path>")
		os.Exit(1)
	}

	// Run the main logic and handle errors
	if err := run(*dictionaryPath, *inputPath); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
