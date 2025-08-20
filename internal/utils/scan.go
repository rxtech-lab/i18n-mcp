package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// ScanPoFiles scans all .po files in the given path and returns a list of file paths
func ScanPoFiles(path string) ([]string, error) {
	var poFiles []string

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file has .po extension
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(filePath), ".po") {
			poFiles = append(poFiles, filePath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return poFiles, nil
}
