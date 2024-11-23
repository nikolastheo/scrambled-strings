package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"scrambled-strings/pkg/matcher"
)

// CLIArgs holds the command-line arguments for validation
type CLIArgs struct {
	DictionaryPath string `validate:"required,file,readable"`
	InputPath      string `validate:"required,file,readable"`
}

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

// fileExists checks if a file exists.
func fileExists(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	_, err := os.Stat(path)
	return err == nil
}

// fileReadable checks if a file is readable.
func fileReadable(fl validator.FieldLevel) bool {
	path := fl.Field().String()
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()
	return true
}

// validateArgs validates the CLI arguments
func validateArgs(dictionaryPath, inputPath string) error {
	validate := validator.New()

	// Register custom validations
	validate.RegisterValidation("file", fileExists)
	validate.RegisterValidation("readable", fileReadable)

	args := CLIArgs{
		DictionaryPath: dictionaryPath,
		InputPath:      inputPath,
	}

	if err := validate.Struct(args); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			if e.Field() == "DictionaryPath" {
				return fmt.Errorf("dictionary file validation failed: %s", e.Tag())
			}
			if e.Field() == "InputPath" {
				return fmt.Errorf("input file validation failed: %s", e.Tag())
			}
		}
	}

	return nil
}

// run processes the dictionary and input files and counts matches.
func run(dictionaryPath, inputPath string) error {
	// Validate arguments
	if err := validateArgs(dictionaryPath, inputPath); err != nil {
		return err
	}

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
