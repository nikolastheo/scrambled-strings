
# Scrambled-Strings

## Overview

The **Scrambled-Strings** project is a command-line tool written in Go that counts how many dictionary words (in their original or scrambled forms) appear as substrings in a set of input strings. The scrambled form of a dictionary word maintains the first and last letters while the middle characters can be rearranged.

---

## Requirements

Before running the project, ensure you have the following installed and up to date:

### Go

You can check if Go is installed and its version by running:

```bash
go version
```

If not installed, visit [Go Installation Guide](https://golang.org/doc/install).

### Docker (Optional)

Docker is required if you wish to run the application inside a container. Verify Docker installation by running:

```bash
docker --version
```

If not installed, visit [Docker Installation Guide](https://docs.docker.com/get-docker/).

---

## Usage

### Build and Run Locally

#### Build the Application

From the project root directory, build the application using:

```bash
go build -o scrambled-strings ./cmd
```

This creates an executable named `scrambled-strings`.

#### Run the Application

The application requires two input files:

1. **Dictionary File**: Each line contains a dictionary word.
2. **Input File**: Each line contains a long string to search within.

Test Data:
```
├── testdata
    ├── dictionary.txt
    ├── input.txt
```

Run the application using the following syntax:

```bash
./scrambled-strings --dictionary <dictionary_file_path> --input <input_file_path> [--log-level <log_level>]
```

- `--dictionary`: Path to the dictionary file (required).
- `--input`: Path to the input file (required).
- `--log-level`: Log verbosity level. Options are `debug`, `info`, `warn`, `error`, `fatal`, `panic`. Default is `info`.

#### Example Usage

```bash
./scrambled-strings --dictionary data/dictionary.txt --input data/input.txt
```

With debug-level logging:

```bash
./scrambled-strings --dictionary data/dictionary.txt --input data/input.txt --log-level debug
```
for help:

```bash
./scrambled-strings --help
```

### Sample Input Files

#### Dictionary File (`dictionary.txt`)

```
axpaj
apxaj
dnrbt
pjxdn
abd
```

#### Input File (`input.txt`)

```
aapxjdnrbtvldptfzbbdbbzxtndrvjblnzjfpvhdhhpxjdnrbt
nothingmatcheshere
```

#### Expected Output

```
Case #1: 4
Case #2: 0
```

---

## Running with Docker

### Build Docker Image

From the project root directory, build the Docker image:

```bash
docker build -t scrambled-strings .
```

### Prepare Input Files

1. Create a `data` directory inside your project folder or anywhere in your system.
2. Place your `dictionary.txt` and `input.txt` files inside the `data` directory.

For example:

```
/data
|--> dictionary.txt
|--> input.txt
```

### Run Docker Container

Run the application using Docker:

```bash
docker run -v <absolute_path_to_data>:/app/data scrambled-strings --dictionary /app/data/dictionary.txt --input /app/data/input.txt
```

#### Example

```bash
docker run -v $(pwd)/data:/app/data scrambled-strings --dictionary /app/data/dictionary.txt --input /app/data/input.txt
```

---

## Testing

Unit tests are included to verify the application's functionality.

### Run Tests

Run all tests using:

```bash
go test ./...
```

---

## Design Concepts

### Overview

1. **Parse Inputs**: Read dictionary and input files.
2. **Canonical Form**: Generate a canonical form for each word in the dictionary to account for scrambled forms.
3. **Matching Logic**: Iterate through input strings to find matches based on canonical forms.
4. **Output Results**: Display results in the format `Case #x: y`.

### Key Modules

#### Canonical Form

**Function**: Ensures the first and last letters are preserved while sorting the middle characters.

#### Matching Logic

- For each input string:
  - Extract substrings matching the length of dictionary words.
  - Check if substrings match the original or canonical forms of dictionary words.
  - Avoid double counting matches for the same word.

---

## Configuration Options

Modify configurations (e.g., logging levels, buffer sizes) using a configuration file.

---

## Future Improvements

1. Optimize matching with algorithms like Aho-Corasick for multi-pattern matching.
2. Implement concurrency to process multiple input strings in parallel.
3. Add extended logging for detailed analytics.

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
