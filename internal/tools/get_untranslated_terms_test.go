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

func TestGetUntranslatedTermsTool(t *testing.T) {
	// Create temporary directory for test files
	tempDir, err := os.MkdirTemp("", "po_untranslated_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a test .po file with mixed translated and untranslated terms
	poContent := `# Test PO file
msgid ""
msgstr ""
"Content-Type: text/plain; charset=UTF-8\n"
"Language: en\n"

msgid "hello"
msgstr "hola"

msgid "goodbye"
msgstr "adiós"

msgid "untranslated1"
msgstr ""

msgid "untranslated2"
msgstr ""

msgid "untranslated3"
msgstr ""

msgid "same_as_key"
msgstr "same_as_key"

msgid "untranslated4"
msgstr ""
`

	poFile := filepath.Join(tempDir, "test.po")
	err = os.WriteFile(poFile, []byte(poContent), 0644)
	require.NoError(t, err)

	// Get the tool and handler
	tool, handler := NewGetUntranslatedTermsTool()

	// Verify tool properties
	assert.Equal(t, "getUntranslatedTerms", tool.Name)
	assert.Contains(t, tool.Description, "untranslated terms")

	// Test with valid file and default limit
	t.Run("Valid File Default Limit", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"file_path": poFile,
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)

		// Parse the result
		var resultData map[string]interface{}
		textContent := getTextContent(t, result)
		err = json.Unmarshal([]byte(textContent), &resultData)
		require.NoError(t, err)

		assert.Equal(t, poFile, resultData["file_path"])
		assert.Equal(t, float64(10), resultData["limit"])

		untranslatedTerms := resultData["untranslated_terms"].(map[string]interface{})
		// Should have untranslated terms (empty msgstr or same as msgid)
		assert.Greater(t, len(untranslatedTerms), 0)
		assert.LessOrEqual(t, len(untranslatedTerms), 10)
	})

	// Test with custom limit
	t.Run("Custom Limit", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"file_path": poFile,
			"limit":     "2",
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)

		var resultData map[string]interface{}
		textContent := getTextContent(t, result)
		err = json.Unmarshal([]byte(textContent), &resultData)
		require.NoError(t, err)

		assert.Equal(t, float64(2), resultData["limit"])
		assert.Equal(t, float64(2), resultData["count"])

		untranslatedTerms := resultData["untranslated_terms"].(map[string]interface{})
		assert.Equal(t, 2, len(untranslatedTerms))
	})

	// Test with non-existent file
	t.Run("Non-existent File", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"file_path": "/non/existent/file.po",
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)
		assert.True(t, result.IsError)
		textContent := getTextContent(t, result)
		assert.Contains(t, textContent, "Error parsing PO file")
	})

	// Test with missing file_path parameter
	t.Run("Missing FilePath Parameter", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{})

		_, err := handler(context.Background(), request)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "file_path parameter is required")
	})

	// Test with invalid limit
	t.Run("Invalid Limit", func(t *testing.T) {
		request := makeRequest(map[string]interface{}{
			"file_path": poFile,
			"limit":     "invalid",
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)
		assert.True(t, result.IsError)
		textContent := getTextContent(t, result)
		assert.Contains(t, textContent, "Invalid limit value")
	})

	// Test with all translated file
	t.Run("All Translated File", func(t *testing.T) {
		allTranslatedContent := `# Test PO file
msgid ""
msgstr ""
"Content-Type: text/plain; charset=UTF-8\n"
"Language: en\n"

msgid "hello"
msgstr "hola"

msgid "goodbye"
msgstr "adiós"
`
		allTranslatedFile := filepath.Join(tempDir, "all_translated.po")
		err = os.WriteFile(allTranslatedFile, []byte(allTranslatedContent), 0644)
		require.NoError(t, err)

		request := makeRequest(map[string]interface{}{
			"file_path": allTranslatedFile,
		})

		result, err := handler(context.Background(), request)
		require.NoError(t, err)

		var resultData map[string]interface{}
		textContent := getTextContent(t, result)
		err = json.Unmarshal([]byte(textContent), &resultData)
		require.NoError(t, err)

		assert.Equal(t, float64(0), resultData["count"])
		untranslatedTerms := resultData["untranslated_terms"].(map[string]interface{})
		assert.Empty(t, untranslatedTerms)
	})
}
