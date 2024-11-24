package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"scrambled-strings/pkg/matcher"
	"github.com/go-playground/validator/v10"
)

// Set up logger
var log = logrus.New()

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

// configureLogLevel sets the log level based on the user's input.
func configureLogLevel(logLevel string) {
	level, err := logrus.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		log.WithError(err).Warnf("Invalid log level: '%s'. Defaulting to 'info'", logLevel)
		level = logrus.InfoLevel
	}
	log.SetLevel(level)
	log.WithField("logLevel", level.String()).Info("Log level configured")
}

// run processes the dictionary and input files and counts matches.
func run(dictionaryPath, inputPath string) error {
	log.WithFields(logrus.Fields{
		"dictionaryPath": dictionaryPath,
		"inputPath":      inputPath,
	}).Info("Starting run function")

	// Validate arguments
	if err := validateArgs(dictionaryPath, inputPath); err != nil {
		log.WithError(err).Error("Validation failed")
		return err
	}

	// Read and preprocess the dictionary
	dictionaryWords, err := readLines(dictionaryPath)
	if err != nil {
		log.WithError(err).Error("Error reading dictionary file")
		return fmt.Errorf("error reading dictionary file: %v", err)
	}
	log.WithField("wordCount", len(dictionaryWords)).Info("Loaded dictionary")

	canonicalDictionary := matcher.PrecomputeCanonicalForms(dictionaryWords)
	log.Info("Precomputed canonical forms for dictionary")

	// Read input strings
	inputStrings, err := readLines(inputPath)
	if err != nil {
		log.WithError(err).Error("Error reading input file")
		return fmt.Errorf("error reading input file: %v", err)
	}
	log.WithField("inputLineCount", len(inputStrings)).Info("Loaded input strings")

	// Process each input string and print results
	for i, input := range inputStrings {
		matchCount := matcher.CountMatches(canonicalDictionary, input)
		log.WithFields(logrus.Fields{
			"case":       i + 1,
			"matchCount": matchCount,
		}).Info("Processed input line")
		fmt.Printf("Case #%d: %d\n", i+1, matchCount)
	}

	log.Info("Run function completed successfully")
	return nil
}

func main() {
	// Configure logger
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Define command-line flags
	dictionaryPath := flag.String("dictionary", "", "Path to the dictionary file")
	inputPath := flag.String("input", "", "Path to the input file")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error, fatal, panic)")
	flag.Parse()

	// Ensure both flags are provided
	if *dictionaryPath == "" || *inputPath == "" {
		log.Error("Missing required flags")
		fmt.Println("Usage: ./scrambled-strings --dictionary <path> --input <path> [--log-level <level>]")
		os.Exit(1)
	}

	// Configure log level
	configureLogLevel(*logLevel)

	// Run the main logic and handle errors
	if err := run(*dictionaryPath, *inputPath); err != nil {
		log.WithError(err).Fatal("Program terminated with an error")
		os.Exit(1)
	}
}
