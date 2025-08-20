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

func TestLookUpTranslationTool(t *testing.T) {
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "po_lookup_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a test .po file with various terms
	poContent := `# Test PO file
msgid ""
msgstr ""
"Content-Type: text/plain; charset=UTF-8\n"
"Language: en\n"

msgid "hello"
msgstr "hola"

msgid "hello_world"
msgstr "hola mundo"

msgid "goodbye"
msgstr "adiós"

msgid "welcome_message"
msgstr "mensaje de bienvenida"

msgid "error_not_found"
msgstr "error no encontrado"

msgid "error_invalid"
msgstr "error inválido"

msgid "button_ok"
msgstr "botón ok"

msgid "button_cancel"
msgstr "botón cancelar"
`

	poFile := filepath.Join(tempDir, "test.po")
	err = os.WriteFile(poFile, []byte(poContent), 0644)
	require.NoError(t, err)

	// Get the tool and handler
	tool, handler := NewLookUpTranslationTool()

	// Verify tool properties
	assert.Equal(t, "lookUpTranslation", tool.Name)
	assert.Contains(t, tool.Description, "Search for a term")

	// Test searching for "hello"
	t.Run("Search Hello", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"file_path":   poFile,
			"search_term": "hello",
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)
		assert.NotNil(t, result)

		var resultData map[string]interface{}
		textContent := getTextContent(t, result)
		err = json.Unmarshal([]byte(textContent), &resultData)
		require.NoError(t, err)

		assert.Equal(t, "hello", resultData["search_term"])
		assert.Equal(t, float64(2), resultData["total_matches"]) // "hello" and "hello_world"

		translations := resultData["translations"].(map[string]interface{})
		assert.Contains(t, translations, "hello")
		assert.Contains(t, translations, "hello_world")
	})

	// Test searching for "error" with pagination
	t.Run("Search Error With Pagination", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"file_path":   poFile,
			"search_term": "error",
			"page_size":   "1",
			"page":        "1",
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)
		assert.NotNil(t, result)

		var resultData map[string]interface{}
		textContent := getTextContent(t, result)
		err = json.Unmarshal([]byte(textContent), &resultData)
		require.NoError(t, err)

		assert.Equal(t, float64(2), resultData["total_matches"]) // "error_not_found" and "error_invalid"
		assert.Equal(t, float64(1), resultData["page"])
		assert.Equal(t, float64(1), resultData["page_size"])

		translations := resultData["translations"].(map[string]interface{})
		assert.Len(t, translations, 1)
	})

	// Test searching for "button" - case insensitive
	t.Run("Search Button Case Insensitive", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"file_path":   poFile,
			"search_term": "BUTTON",
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)
		assert.NotNil(t, result)

		var resultData map[string]interface{}
		textContent := getTextContent(t, result)
		err = json.Unmarshal([]byte(textContent), &resultData)
		require.NoError(t, err)

		assert.Equal(t, float64(2), resultData["total_matches"]) // "button_ok" and "button_cancel"
		translations := resultData["translations"].(map[string]interface{})
		assert.Contains(t, translations, "button_ok")
		assert.Contains(t, translations, "button_cancel")
	})

	// Test with non-existent search term
	t.Run("Non-existent Search Term", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"file_path":   poFile,
			"search_term": "nonexistent",
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)
		assert.NotNil(t, result)

		var resultData map[string]interface{}
		textContent := getTextContent(t, result)
		err = json.Unmarshal([]byte(textContent), &resultData)
		require.NoError(t, err)

		assert.Equal(t, float64(0), resultData["total_matches"])
		translations := resultData["translations"].(map[string]interface{})
		assert.Empty(t, translations)
	})

	// Test pagination - page 2
	t.Run("Pagination Page 2", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"file_path":   poFile,
			"search_term": "error",
			"page_size":   "1",
			"page":        "2",
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)
		assert.NotNil(t, result)

		var resultData map[string]interface{}
		textContent := getTextContent(t, result)
		err = json.Unmarshal([]byte(textContent), &resultData)
		require.NoError(t, err)

		assert.Equal(t, float64(2), resultData["page"])
		translations := resultData["translations"].(map[string]interface{})
		assert.Len(t, translations, 1)
	})

	// Test with non-existent file
	t.Run("Non-existent File", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"file_path":   "/non/existent/file.po",
			"search_term": "hello",
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
			"search_term": "hello",
		})

		_, err := handler(context.Background(), request)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "file_path parameter is required")
	})

	t.Run("Missing SearchTerm Parameter", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"file_path": poFile,
		})

		_, err := handler(context.Background(), request)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "search_term parameter is required")
	})

	// Test with invalid pagination parameters
	t.Run("Invalid Page Size", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"file_path":   poFile,
			"search_term": "hello",
			"page_size":   "invalid",
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.IsError)
		textContent := getTextContent(t, result)
		assert.Contains(t, textContent, "Invalid page_size value")
	})

	t.Run("Invalid Page Number", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"file_path":   poFile,
			"search_term": "hello",
			"page":        "invalid",
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.IsError)
		textContent := getTextContent(t, result)
		assert.Contains(t, textContent, "Invalid page value")
	})
}
