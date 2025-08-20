package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListAllPoFilesTool(t *testing.T) {
	// Create temporary directory structure with .po files
	tempDir, err := os.MkdirTemp("", "po_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create subdirectories and .po files
	subDir1 := filepath.Join(tempDir, "locale", "en")
	subDir2 := filepath.Join(tempDir, "locale", "es")
	err = os.MkdirAll(subDir1, 0755)
	require.NoError(t, err)
	err = os.MkdirAll(subDir2, 0755)
	require.NoError(t, err)

	// Create test .po files with language headers
	poFiles := []struct {
		path    string
		content string
	}{
		{filepath.Join(subDir1, "messages.po"), `msgid ""
msgstr ""
"Language: en\n"

msgid "Hello"
msgstr "Hello"`},
		{filepath.Join(subDir1, "errors.po"), `msgid ""
msgstr ""
"Language: en\n"

msgid "Error"
msgstr "Error"`},
		{filepath.Join(subDir2, "messages.po"), `msgid ""
msgstr ""
"Language: es\n"

msgid "Hello"
msgstr "Hola"`},
		{filepath.Join(tempDir, "main.po"), `msgid ""
msgstr ""
"Language: fr\n"

msgid "Main"
msgstr "Principal"`},
	}

	for _, file := range poFiles {
		err = os.WriteFile(file.path, []byte(file.content), 0644)
		require.NoError(t, err)
	}

	// Create a non-.po file that should be ignored
	err = os.WriteFile(filepath.Join(tempDir, "readme.txt"), []byte("Not a PO file"), 0644)
	require.NoError(t, err)

	// Get the tool and handler
	tool, handler := NewListAllPoFilesTool()

	// Verify tool properties
	assert.Equal(t, "listAllPoFiles", tool.Name)
	assert.Contains(t, tool.Description, "List all .po files")

	// Test with valid directory
	t.Run("Valid Directory", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"directory": tempDir,
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)

		// Parse the result
		var resultData map[string]interface{}
		textContent := getTextContent(t, result)
		err = json.Unmarshal([]byte(textContent), &resultData)
		require.NoError(t, err)

		assert.Equal(t, tempDir, resultData["directory"])
		assert.Equal(t, float64(4), resultData["count"])

		files := resultData["files"].([]interface{})
		assert.Len(t, files, 4)

		// Check that files have language information
		for _, f := range files {
			fileInfo := f.(map[string]interface{})
			assert.Contains(t, fileInfo, "path")
			assert.Contains(t, fileInfo, "language")
			// Check language values
			lang := fileInfo["language"].(string)
			assert.Contains(t, []string{"en", "es", "fr"}, lang)
		}
	})

	// Test with non-existent directory
	t.Run("Non-existent Directory", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"directory": "/non/existent/path",
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)
		assert.True(t, result.IsError)
		textContent := getTextContent(t, result)
		assert.Contains(t, textContent, "Error scanning")
	})

	// Test with missing directory parameter
	t.Run("Missing Directory Parameter", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{})

		_, err := handler(context.Background(), request)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "directory parameter is required")
	})

	// Test with empty directory (no .po files)
	t.Run("Empty Directory", func(t *testing.T) {
		emptyDir, err := os.MkdirTemp("", "empty_test")
		require.NoError(t, err)
		defer os.RemoveAll(emptyDir)

		request := makeRequest(map[string]interface{}{
			"directory": emptyDir,
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)

		var resultData map[string]interface{}
		textContent := getTextContent(t, result)
		err = json.Unmarshal([]byte(textContent), &resultData)
		require.NoError(t, err)

		assert.Equal(t, float64(0), resultData["count"])
		files, ok := resultData["files"].([]interface{})
		if ok {
			assert.Len(t, files, 0)
		} else {
			// Files could be nil for empty directory
			assert.Nil(t, resultData["files"])
		}
	})
}
