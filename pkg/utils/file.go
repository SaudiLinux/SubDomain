package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadWordlist loads a wordlist from a file
func LoadWordlist(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open wordlist file: %v", err)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if word != "" && !strings.HasPrefix(word, "#") {
			words = append(words, word)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading wordlist file: %v", err)
	}

	return words, nil
}

// LoadTargets loads targets from a file or returns a single target
func LoadTargets(target string) ([]string, error) {
	// Check if the target is a file
	if _, err := os.Stat(target); err == nil {
		return LoadWordlist(target)
	}

	// Return the target as a single item
	return []string{target}, nil
}

// EnsureDirectory ensures that a directory exists
func EnsureDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// SaveToFile saves content to a file
func SaveToFile(path string, content string) error {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := EnsureDirectory(dir); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Create or truncate the file
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}

// AppendToFile appends content to a file
func AppendToFile(path string, content string) error {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := EnsureDirectory(dir); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Open the file in append mode
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Write content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}

// CreateOutputFile creates an output file with a header
func CreateOutputFile(path string, header string) (*os.File, error) {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := EnsureDirectory(dir); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	// Create or truncate the file
	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err)
	}

	// Write header to the file
	if header != "" {
		_, err = file.WriteString(header + "\n")
		if err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to write header to file: %v", err)
		}
	}

	return file, nil
}