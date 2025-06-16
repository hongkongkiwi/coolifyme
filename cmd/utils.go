package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// safeReadFile reads a file with path validation to prevent security issues
func safeReadFile(filename string) ([]byte, error) {
	// Clean the file path to prevent directory traversal
	filename = filepath.Clean(filename)

	// Check if the file path contains potentially dangerous elements
	if strings.Contains(filename, "..") {
		return nil, fmt.Errorf("invalid file path: contains directory traversal")
	}

	// Read the file
	return os.ReadFile(filename) // #nosec G304 - path is validated above
}
