package tools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTranslateTool(t *testing.T) {
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "po_translate_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Get the tool and handler
	tool, handler := NewTranslateTool()

	// Verify tool properties
	assert.Equal(t, "translate", tool.Name)
	assert.Contains(t, tool.Description, "Translate terms")

	// Test translating multiple terms
	t.Run("Translate Multiple Terms", func(t *testing.T) {
		// Create initial PO file
		poContent := `# Test PO file
msgid ""
msgstr ""
"Content-Type: text/plain; charset=UTF-8\n"
"Language: en\n"

msgid "hello"
msgstr ""

msgid "goodbye"
msgstr ""

msgid "welcome"
msgstr ""

msgid "thank_you"
msgstr "gracias"
`
		poFile := filepath.Join(tempDir, "test_translate.po")
		err = os.WriteFile(poFile, []byte(poContent), 0644)
		require.NoError(t, err)

		// Prepare translations
		translations := map[string]string{
			"hello":   "hola",
			"goodbye": "adiós",
			"welcome": "bienvenido",
		}

		translationsJSON, err := json.Marshal(translations)
		require.NoError(t, err)

		request := makeRequest(map[string]interface{}{
			"file_path":    poFile,
			"translations": string(translationsJSON),
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.IsError)

		// Parse the result
		var resultData map[string]interface{}
		textContent := getTextContent(t, result)
		err = json.Unmarshal([]byte(textContent), &resultData)
		require.NoError(t, err)

		assert.Equal(t, float64(3), resultData["translated_count"])
		assert.Contains(t, resultData["message"].(string), "Successfully translated 3 terms")

		// Verify the file was updated
		updatedContent, err := os.ReadFile(poFile)
		require.NoError(t, err)
		updatedStr := string(updatedContent)

		// Check that translations were applied
		assert.Contains(t, updatedStr, "hola")
		assert.Contains(t, updatedStr, "adiós")
		assert.Contains(t, updatedStr, "bienvenido")
		// Original translation should remain
		assert.Contains(t, updatedStr, "gracias")
	})

	// Test with empty translations
	t.Run("Empty Translations", func(t *testing.T) {
		poContent := `# Test PO file
msgid ""
msgstr ""
"Content-Type: text/plain; charset=UTF-8\n"
"Language: en\n"

msgid "test"
msgstr ""
`
		poFile := filepath.Join(tempDir, "test_empty.po")
		err = os.WriteFile(poFile, []byte(poContent), 0644)
		require.NoError(t, err)

		request := makeRequest(map[string]interface{}{
			"file_path":    poFile,
			"translations": "{}",
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.IsError)

		var resultData map[string]interface{}
		err = json.Unmarshal([]byte(getTextContent(t, result)), &resultData)
		require.NoError(t, err)

		assert.Equal(t, float64(0), resultData["translated_count"])
	})

	// Test with invalid JSON
	t.Run("Invalid JSON", func(t *testing.T) {
		poFile := filepath.Join(tempDir, "test_invalid.po")
		err = os.WriteFile(poFile, []byte("# Test"), 0644)
		require.NoError(t, err)

		request := makeRequest(map[string]interface{}{
			"file_path":    poFile,
			"translations": "invalid json",
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.IsError)
		textContent := getTextContent(t, result)
		assert.Contains(t, textContent, "Invalid translations JSON")
	})

	// Test with non-existent file
	t.Run("Non-existent File", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"file_path":    "/non/existent/file.po",
			"translations": "{}",
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.IsError)
		textContent := getTextContent(t, result)
		assert.Contains(t, textContent, "Error parsing PO file")
	})

	// Test with missing parameters
	t.Run("Missing FilePath Parameter", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"translations": "{}",
		})

		_, err := handler(context.Background(), request)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "file_path parameter is required")
	})

	t.Run("Missing Translations Parameter", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"file_path": "/some/file.po",
		})

		_, err := handler(context.Background(), request)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "translations parameter is required")
	})

	// Test updating existing translations
	t.Run("Update Existing Translations", func(t *testing.T) {
		poContent := `# Test PO file
msgid ""
msgstr ""
"Content-Type: text/plain; charset=UTF-8\n"
"Language: en\n"

msgid "hello"
msgstr "old_translation"

msgid "goodbye"
msgstr "old_goodbye"
`
		poFile := filepath.Join(tempDir, "test_update.po")
		err = os.WriteFile(poFile, []byte(poContent), 0644)
		require.NoError(t, err)

		translations := map[string]string{
			"hello":   "hola",
			"goodbye": "adiós",
		}

		translationsJSON, err := json.Marshal(translations)
		require.NoError(t, err)

		request := makeRequest(map[string]interface{}{
			"file_path":    poFile,
			"translations": string(translationsJSON),
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.IsError)

		// Verify the file was updated with new translations
		updatedContent, err := os.ReadFile(poFile)
		require.NoError(t, err)
		updatedStr := string(updatedContent)

		assert.Contains(t, updatedStr, "hola")
		assert.Contains(t, updatedStr, "adiós")
		assert.NotContains(t, updatedStr, "old_translation")
		assert.NotContains(t, updatedStr, "old_goodbye")
	})

	// Test file permissions (read-only file)
	t.Run("Read-only File", func(t *testing.T) {
		poContent := `# Test PO file
msgid ""
msgstr ""
"Content-Type: text/plain; charset=UTF-8\n"
"Language: en\n"

msgid "test"
msgstr ""
`
		poFile := filepath.Join(tempDir, "test_readonly.po")
		err = os.WriteFile(poFile, []byte(poContent), 0444) // Read-only
		require.NoError(t, err)

		// Try to make it writable for cleanup
		defer os.Chmod(poFile, 0644)

		translations := map[string]string{
			"test": "prueba",
		}

		translationsJSON, err := json.Marshal(translations)
		require.NoError(t, err)

		request := makeRequest(map[string]interface{}{
			"file_path":    poFile,
			"translations": string(translationsJSON),
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.IsError)
		textContent := getTextContent(t, result)
		assert.Contains(t, strings.ToLower(textContent), "error writing")
	})
}
