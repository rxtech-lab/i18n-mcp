package utils

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestScanPoFilesWithInfo(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "po-test-with-info")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test PO files with different languages
	testFiles := []struct {
		name    string
		content string
		lang    string
	}{
		{
			name: "en.po",
			content: `msgid ""
msgstr ""
"Language: en\n"
"Content-Type: text/plain; charset=UTF-8\n"

msgid "Hello"
msgstr "Hello"`,
			lang: "en",
		},
		{
			name: "fr.po",
			content: `msgid ""
msgstr ""
"Language: fr\n"
"Content-Type: text/plain; charset=UTF-8\n"

msgid "Hello"
msgstr "Bonjour"`,
			lang: "fr",
		},
		{
			name: "zh-HK.po",
			content: `msgid ""
msgstr ""
"Language: zh-HK\n"
"Content-Type: text/plain; charset=UTF-8\n"

msgid "Hello"
msgstr "你好"`,
			lang: "zh-HK",
		},
	}

	// Create the test files
	for _, tf := range testFiles {
		filePath := filepath.Join(tempDir, tf.name)
		err := os.WriteFile(filePath, []byte(tf.content), 0644)
		require.NoError(t, err)
	}

	// Create a non-PO file that should be ignored
	txtFile := filepath.Join(tempDir, "readme.txt")
	err = os.WriteFile(txtFile, []byte("This is not a PO file"), 0644)
	require.NoError(t, err)

	// Create a subdirectory with a PO file
	subDir := filepath.Join(tempDir, "subdir")
	err = os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	subPo := filepath.Join(subDir, "es.po")
	err = os.WriteFile(subPo, []byte(`msgid ""
msgstr ""
"Language: es\n"

msgid "Hello"
msgstr "Hola"`), 0644)
	require.NoError(t, err)

	t.Run("ScanPoFilesWithInfo returns correct information", func(t *testing.T) {
		poFiles, err := ScanPoFilesWithInfo(tempDir)
		require.NoError(t, err)

		// Should find 4 PO files (3 in root + 1 in subdir)
		assert.Len(t, poFiles, 4)

		// Create a map for easier testing
		fileMap := make(map[string]string)
		for _, pf := range poFiles {
			fileMap[filepath.Base(pf.Path)] = pf.Language
		}

		// Check that languages are correctly extracted
		assert.Equal(t, "en", fileMap["en.po"])
		assert.Equal(t, "fr", fileMap["fr.po"])
		assert.Equal(t, "zh-HK", fileMap["zh-HK.po"])
		assert.Equal(t, "es", fileMap["es.po"])
	})

	t.Run("Handles PO file with parsing error gracefully", func(t *testing.T) {
		invalidDir, err := os.MkdirTemp("", "invalid-po-test")
		require.NoError(t, err)
		defer os.RemoveAll(invalidDir)

		// Create an invalid PO file
		invalidPo := filepath.Join(invalidDir, "invalid.po")
		err = os.WriteFile(invalidPo, []byte("This is not valid PO content"), 0644)
		require.NoError(t, err)

		poFiles, err := ScanPoFilesWithInfo(invalidDir)
		require.NoError(t, err)
		assert.Len(t, poFiles, 1)

		// Language should be empty for invalid PO file
		assert.Equal(t, "", poFiles[0].Language)
		assert.Contains(t, poFiles[0].Path, "invalid.po")
	})
}
