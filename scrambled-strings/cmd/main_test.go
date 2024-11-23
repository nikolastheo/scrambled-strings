package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// Helper function to capture stdout during function execution
func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// TestRun validates the main functionality of the `run` function
func TestRun(t *testing.T) {
	// Create temporary test files
	dictionaryPath := "test_dictionary.txt"
	inputPath := "test_input.txt"
	defer os.Remove(dictionaryPath) // Clean up after test
	defer os.Remove(inputPath)

	// Write test data to temporary files
	os.WriteFile(dictionaryPath, []byte("axpaj\napxaj\ndnrbt\npjxdn\nabd\n"), 0644)
	os.WriteFile(inputPath, []byte("aapxjdnrbtvldptfzbbdbbzxtndrvjblnzjfpvhdhhpxjdnrbt\nnothingmatcheshere\n"), 0644)

	output := captureOutput(func() {
		err := run(dictionaryPath, inputPath)
		if err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		}
	})

	expectedOutput := "Case #1: 4\nCase #2: 0\n"
	if strings.TrimSpace(output) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected output:\n%s\nBut got:\n%s", expectedOutput, output)
	}
}

// TestValidateArgs validates the argument validation logic
func TestValidateArgs(t *testing.T) {
	// Create temporary test files
	dictionaryPath := "test_dictionary.txt"
	inputPath := "test_input.txt"
	defer os.Remove(dictionaryPath)
	defer os.Remove(inputPath)

	os.WriteFile(dictionaryPath, []byte("word1\nword2\n"), 0644)
	os.WriteFile(inputPath, []byte("input line\n"), 0644)

	tests := []struct {
		dictionaryPath string
		inputPath      string
		expectError    bool
	}{
		{dictionaryPath, inputPath, false},
		{"nonexistent.txt", inputPath, true},
		{dictionaryPath, "nonexistent-input.txt", true},
		{"nonexistent.txt", "nonexistent-input.txt", true},
	}

	for _, test := range tests {
		err := validateArgs(test.dictionaryPath, test.inputPath)
		if test.expectError && err == nil {
			t.Errorf("Expected an error but got none for dictionary: %s, input: %s", test.dictionaryPath, test.inputPath)
		}
		if !test.expectError && err != nil {
			t.Errorf("Did not expect an error but got: %v for dictionary: %s, input: %s", err, test.dictionaryPath, test.inputPath)
		}
	}
}

// TestRunInvalidDictionary checks behavior when the dictionary file is invalid
func TestRunInvalidDictionary(t *testing.T) {
	err := run("nonexistent_dictionary.txt", "test_input.txt")
	if err == nil || !strings.Contains(err.Error(), "dictionary file validation failed") {
		t.Errorf("Expected validation error for invalid dictionary file but got: %v", err)
	}
}

// TestRunInvalidInput checks behavior when the input file is invalid
func TestRunInvalidInput(t *testing.T) {
	// Create a valid dictionary file
	dictionaryPath := "test_dictionary.txt"
	defer os.Remove(dictionaryPath)
	os.WriteFile(dictionaryPath, []byte("axpaj\napxaj\ndnrbt\npjxdn\nabd\n"), 0644)

	err := run(dictionaryPath, "nonexistent_input.txt")
	if err == nil || !strings.Contains(err.Error(), "input file validation failed") {
		t.Errorf("Expected validation error for invalid input file but got: %v", err)
	}
}

// TestRunEmptyDictionary tests behavior with an empty dictionary
func TestRunEmptyDictionary(t *testing.T) {
	// Create temporary test files
	dictionaryPath := "empty_dictionary.txt"
	inputPath := "test_input.txt"
	defer os.Remove(dictionaryPath)
	defer os.Remove(inputPath)

	os.WriteFile(dictionaryPath, []byte(""), 0644)
	os.WriteFile(inputPath, []byte("aapxjdnrbtvldptfzbbdbbzxtndrvjblnzjfpvhdhhpxjdnrbt\n"), 0644)

	output := captureOutput(func() {
		err := run(dictionaryPath, inputPath)
		if err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		}
	})

	expectedOutput := "Case #1: 0\n"
	if strings.TrimSpace(output) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected output:\n%s\nBut got:\n%s", expectedOutput, output)
	}
}

// TestRunEmptyInput tests behavior with an empty input file
func TestRunEmptyInput(t *testing.T) {
	dictionaryPath := "test_dictionary.txt"
	inputPath := "empty_input.txt"
	defer os.Remove(dictionaryPath)
	defer os.Remove(inputPath)

	os.WriteFile(dictionaryPath, []byte("axpaj\napxaj\ndnrbt\npjxdn\nabd\n"), 0644)
	os.WriteFile(inputPath, []byte(""), 0644)

	output := captureOutput(func() {
		err := run(dictionaryPath, inputPath)
		if err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		}
	})

	expectedOutput := ""
	if strings.TrimSpace(output) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected no output but got:\n%s", output)
	}
}
