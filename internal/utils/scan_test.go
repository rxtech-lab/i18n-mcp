package utils

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestScanPoFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scan_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directory structure
	testFiles := []string{
		"messages.po",
		"translations.po",
		"subdir1/app.po",
		"subdir1/nested/deep.po",
		"subdir2/locale.po",
		"subdir2/nested/another.po",
		"notpo.txt",
		"readme.md",
		"subdir3/test.txt",
	}

	// Create directories and files
	for _, file := range testFiles {
		fullPath := filepath.Join(tempDir, file)
		dir := filepath.Dir(fullPath)

		// Create directory if it doesn't exist
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}

		// Create file
		if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	// Run ScanPoFiles
	poFiles, err := ScanPoFiles(tempDir)
	if err != nil {
		t.Fatalf("ScanPoFiles failed: %v", err)
	}

	// Expected .po files
	expectedFiles := []string{
		filepath.Join(tempDir, "messages.po"),
		filepath.Join(tempDir, "translations.po"),
		filepath.Join(tempDir, "subdir1/app.po"),
		filepath.Join(tempDir, "subdir1/nested/deep.po"),
		filepath.Join(tempDir, "subdir2/locale.po"),
		filepath.Join(tempDir, "subdir2/nested/another.po"),
	}

	// Sort both slices for comparison
	sort.Strings(poFiles)
	sort.Strings(expectedFiles)

	// Check if the number of files matches
	if len(poFiles) != len(expectedFiles) {
		t.Errorf("Expected %d .po files, got %d", len(expectedFiles), len(poFiles))
		t.Errorf("Found files: %v", poFiles)
		t.Errorf("Expected files: %v", expectedFiles)
		return
	}

	// Check each file
	for i, expectedFile := range expectedFiles {
		if poFiles[i] != expectedFile {
			t.Errorf("Expected file %s, got %s", expectedFile, poFiles[i])
		}
	}
}

func TestScanPoFiles_EmptyDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scan_test_empty")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Run ScanPoFiles on empty directory
	poFiles, err := ScanPoFiles(tempDir)
	if err != nil {
		t.Fatalf("ScanPoFiles failed: %v", err)
	}

	if len(poFiles) != 0 {
		t.Errorf("Expected 0 .po files in empty directory, got %d", len(poFiles))
	}
}

func TestScanPoFiles_NonExistentDirectory(t *testing.T) {
	// Try to scan a non-existent directory
	_, err := ScanPoFiles("/path/that/does/not/exist")
	if err == nil {
		t.Error("Expected error when scanning non-existent directory, got nil")
	}
}

func TestScanPoFiles_CaseInsensitive(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "scan_test_case")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files with different case extensions
	testFiles := []string{
		"file1.po",
		"file2.Po",
		"file3.PO",
		"file4.pO",
	}

	for _, file := range testFiles {
		fullPath := filepath.Join(tempDir, file)
		if err := os.WriteFile(fullPath, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}
	}

	// Run ScanPoFiles
	poFiles, err := ScanPoFiles(tempDir)
	if err != nil {
		t.Fatalf("ScanPoFiles failed: %v", err)
	}

	// All files should be found regardless of case
	if len(poFiles) != len(testFiles) {
		t.Errorf("Expected %d .po files (case-insensitive), got %d", len(testFiles), len(poFiles))
		t.Errorf("Found files: %v", poFiles)
	}
}
