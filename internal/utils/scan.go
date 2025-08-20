package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// PoFileInfo contains information about a PO file
type PoFileInfo struct {
	Path     string `json:"path"`
	Language string `json:"language"`
}

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

// ScanPoFilesWithInfo scans all .po files in the given path and returns detailed information
func ScanPoFilesWithInfo(path string) ([]PoFileInfo, error) {
	var poFilesInfo []PoFileInfo

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file has .po extension
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(filePath), ".po") {
			// Try to parse the PO file to get language info
			po, err := ParsePoFile(filePath)
			language := ""
			if err == nil {
				language = po.Language
			}

			poFilesInfo = append(poFilesInfo, PoFileInfo{
				Path:     filePath,
				Language: language,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return poFilesInfo, nil
}
