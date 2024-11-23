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
	oldStdout := os.Stdout              // Save the current stdout
	r, w, _ := os.Pipe()               // Create a pipe to capture output
	os.Stdout = w                      // Redirect stdout to the pipe
	f()                                // Execute the function
	w.Close()                          // Close the write end of the pipe
	os.Stdout = oldStdout              // Restore original stdout
	var buf bytes.Buffer
	io.Copy(&buf, r)                   // Read the captured output into a buffer
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

	// Capture the output of the `run` function
	output := captureOutput(func() {
		err := run(dictionaryPath, inputPath)
		if err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		}
	})

	// Expected output
	expectedOutput := "Case #1: 4\nCase #2: 0\n"

	// Validate the output
	if strings.TrimSpace(output) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected output:\n%s\nBut got:\n%s", expectedOutput, output)
	}
}

// TestRunInvalidDictionary checks behavior when the dictionary file is invalid
func TestRunInvalidDictionary(t *testing.T) {
	// Pass a non-existent dictionary file
	err := run("nonexistent_dictionary.txt", "test_input.txt")
	if err == nil || !strings.Contains(err.Error(), "error reading dictionary file") {
		t.Errorf("Expected error for invalid dictionary file but got: %v", err)
	}
}

// TestRunInvalidInput checks behavior when the input file is invalid
func TestRunInvalidInput(t *testing.T) {
	// Create a valid dictionary file
	dictionaryPath := "test_dictionary.txt"
	defer os.Remove(dictionaryPath) // Clean up after test
	os.WriteFile(dictionaryPath, []byte("axpaj\napxaj\ndnrbt\npjxdn\nabd\n"), 0644)

	// Pass a non-existent input file
	err := run(dictionaryPath, "nonexistent_input.txt")
	if err == nil || !strings.Contains(err.Error(), "error reading input file") {
		t.Errorf("Expected error for invalid input file but got: %v", err)
	}
}

// TestRunEmptyDictionary tests behavior with an empty dictionary
func TestRunEmptyDictionary(t *testing.T) {
	// Create temporary test files
	dictionaryPath := "empty_dictionary.txt"
	inputPath := "test_input.txt"
	defer os.Remove(dictionaryPath) // Clean up after test
	defer os.Remove(inputPath)

	// Write empty dictionary file
	os.WriteFile(dictionaryPath, []byte(""), 0644)
	os.WriteFile(inputPath, []byte("aapxjdnrbtvldptfzbbdbbzxtndrvjblnzjfpvhdhhpxjdnrbt\n"), 0644)

	// Capture the output
	output := captureOutput(func() {
		err := run(dictionaryPath, inputPath)
		if err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		}
	})

	// Expected output
	expectedOutput := "Case #1: 0\n"

	// Validate the output
	if strings.TrimSpace(output) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected output:\n%s\nBut got:\n%s", expectedOutput, output)
	}
}

// TestRunEmptyInput tests behavior with an empty input file
func TestRunEmptyInput(t *testing.T) {
	// Create temporary test files
	dictionaryPath := "test_dictionary.txt"
	inputPath := "empty_input.txt"
	defer os.Remove(dictionaryPath) // Clean up after test
	defer os.Remove(inputPath)

	// Write test data
	os.WriteFile(dictionaryPath, []byte("axpaj\napxaj\ndnrbt\npjxdn\nabd\n"), 0644)
	os.WriteFile(inputPath, []byte(""), 0644)

	// Capture the output
	output := captureOutput(func() {
		err := run(dictionaryPath, inputPath)
		if err != nil {
			t.Fatalf("Expected no error but got: %v", err)
		}
	})

	// Since the input is empty, no matches should be found
	expectedOutput := ""

	// Validate the output
	if strings.TrimSpace(output) != strings.TrimSpace(expectedOutput) {
		t.Errorf("Expected no output but got:\n%s", output)
	}
}
