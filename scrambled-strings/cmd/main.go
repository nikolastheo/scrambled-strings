package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"scrambled-strings/pkg/matcher"
)

// Set up logger
var log = logrus.New()

// CLIArgs holds the command-line arguments for validation
type CLIArgs struct {
	DictionaryPath string `validate:"required,file,readable"`
	InputPath      string `validate:"required,file,readable"`
}

// Custom usage function to display help information
func usage() {
	fmt.Println("Usage:")
	fmt.Println("  scrambled-strings --dictionary <path> --input <path> [--log-level <level>]")
	fmt.Println("\nOptions:")
	fmt.Println("  --dictionary <path>  Path to the dictionary file (required).")
	fmt.Println("  --input <path>       Path to the input file (required).")
	fmt.Println("  --log-level <level>  Log verbosity level: debug, info, warn, error, fatal, panic. Default: 'info'.")
	fmt.Println("\nDescription:")
	fmt.Println("  Counts dictionary words (original or scrambled) appearing as substrings in input strings.")
	fmt.Println("\nExamples:")
	fmt.Println("  scrambled-strings --dictionary /path/to/dictionary.txt --input /path/to/input.txt")
	fmt.Println("  scrambled-strings --dictionary /path/to/dictionary.txt --input /path/to/input.txt --log-level debug")
}

// readLines reads a file and returns its lines as a slice of strings.
func readLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file '%s': %v", filePath, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file '%s': %v", filePath, err)
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
			switch e.Field() {
			case "DictionaryPath":
				return fmt.Errorf("dictionary file validation failed: %s", e.Tag())
			case "InputPath":
				return fmt.Errorf("input file validation failed: %s", e.Tag())
			}
		}
	}

	return nil
}

// configureLogLevel sets the log level based on the user's input.
func configureLogLevel(level string) {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	switch strings.ToLower(level) {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	default:
		log.Warnf("Invalid log level '%s'. Defaulting to 'info'.", level)
		log.SetLevel(logrus.InfoLevel)
	}
}

// run processes the dictionary and input files and counts matches.
func run(dictionaryPath, inputPath string) error {
	log.WithFields(logrus.Fields{
		"dictionaryPath": dictionaryPath,
		"inputPath":      inputPath,
	}).Info("Starting run")

	// Validate arguments
	if err := validateArgs(dictionaryPath, inputPath); err != nil {
		log.WithError(err).Error("Argument validation failed")
		return err
	}

	// Load dictionary words
	dictionaryWords, err := readLines(dictionaryPath)
	if err != nil {
		log.WithError(err).Error("Failed to read dictionary file")
		return err
	}
	log.WithField("wordCount", len(dictionaryWords)).Info("Loaded dictionary words")

	// Precompute canonical forms
	canonicalDictionary := matcher.PrecomputeCanonicalForms(dictionaryWords)
	log.Info("Precomputed canonical forms for dictionary words")

	// Load input strings
	inputStrings, err := readLines(inputPath)
	if err != nil {
		log.WithError(err).Error("Failed to read input file")
		return err
	}
	log.WithField("inputLineCount", len(inputStrings)).Info("Loaded input strings")

	// Process input strings
	for i, input := range inputStrings {
		matchCount := matcher.CountMatches(canonicalDictionary, input)
		log.WithFields(logrus.Fields{
			"case":       i + 1,
			"matchCount": matchCount,
		}).Info("Processed input string")
		fmt.Printf("Case #%d: %d\n", i+1, matchCount)
	}

	log.Info("Run completed successfully")
	return nil
}

func main() {
	// Define command-line flags
	dictionaryPath := flag.String("dictionary", "", "Path to the dictionary file")
	inputPath := flag.String("input", "", "Path to the input file")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error, fatal, panic). Defaults to 'info'.")
	flag.Usage = usage
	flag.Parse()

	// Ensure mandatory flags are provided
	if *dictionaryPath == "" || *inputPath == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Configure log level
	configureLogLevel(*logLevel)

	// Run the main logic and handle errors
	if err := run(*dictionaryPath, *inputPath); err != nil {
		log.WithError(err).Fatal("Execution terminated with an error")
		os.Exit(1)
	}
}
