package service

import (
	"strings"
	"testing"

	"github.com/leonelquinteros/gotext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestPo() *gotext.Po {
	po := gotext.NewPo()
	
	// Parse a sample PO content
	poContent := `msgid ""
msgstr ""
"Content-Type: text/plain; charset=UTF-8\n"
"Language: en\n"
"Plural-Forms: nplurals=2; plural=(n != 1);\n"

msgid "Hello"
msgstr "Hello"

msgid "World"
msgstr "World"

msgid "Untranslated"
msgstr ""

msgid "Same as key"
msgstr "Same as key"

msgid "Different translation"
msgstr "Translated value"

msgid "Another untranslated"
msgstr ""
`
	po.Parse([]byte(poContent))
	return po
}

func TestNewPoService(t *testing.T) {
	po := createTestPo()
	service := NewPoService(*po)
	
	assert.NotNil(t, service)
	assert.NotNil(t, service.poFile)
}

func TestListAllUntranslated(t *testing.T) {
	po := createTestPo()
	service := NewPoService(*po)
	
	t.Run("List all untranslated without limit", func(t *testing.T) {
		untranslated := service.ListAllUntranslated(0)
		
		// Should include "Untranslated" and "Another untranslated" 
		assert.Contains(t, untranslated, "Untranslated")
		assert.Contains(t, untranslated, "Another untranslated")
		
		// Should not include translated items
		assert.NotContains(t, untranslated, "Different translation")
	})
	
	t.Run("List untranslated with limit", func(t *testing.T) {
		untranslated := service.ListAllUntranslated(1)
		
		// Should only have 1 item due to limit
		assert.Len(t, untranslated, 1)
	})
	
	t.Run("List untranslated with large limit", func(t *testing.T) {
		untranslated := service.ListAllUntranslated(100)
		
		// Should include all untranslated items
		assert.Contains(t, untranslated, "Untranslated")
		assert.Contains(t, untranslated, "Another untranslated")
	})
}

func TestTranslate(t *testing.T) {
	po := createTestPo()
	service := NewPoService(*po)
	
	t.Run("Translate existing key", func(t *testing.T) {
		// Translate an untranslated key
		service.Translate("Untranslated", "Now translated")
		
		// Verify it's translated
		assert.Equal(t, "Now translated", service.poFile.Get("Untranslated"))
	})
	
	t.Run("Add new translation", func(t *testing.T) {
		// Add a completely new key
		service.Translate("New Key", "New Value")
		
		// Verify it's added
		assert.Equal(t, "New Value", service.poFile.Get("New Key"))
	})
	
	t.Run("Update existing translation", func(t *testing.T) {
		// Update an already translated key
		service.Translate("Different translation", "Updated value")
		
		// Verify it's updated
		assert.Equal(t, "Updated value", service.poFile.Get("Different translation"))
	})
}

func TestList(t *testing.T) {
	po := createTestPo()
	service := NewPoService(*po)
	
	t.Run("List with skip=0 and take=2", func(t *testing.T) {
		list := service.List(0, 2)
		
		// Should have exactly 2 items
		assert.Len(t, list, 2)
	})
	
	t.Run("List with skip=1 and take=2", func(t *testing.T) {
		list := service.List(1, 2)
		
		// Should have exactly 2 items
		assert.Len(t, list, 2)
	})
	
	t.Run("List with skip exceeding total", func(t *testing.T) {
		list := service.List(100, 2)
		
		// Should have 0 items
		assert.Len(t, list, 0)
	})
	
	t.Run("List with take=0", func(t *testing.T) {
		list := service.List(0, 0)
		
		// Should have 0 items
		assert.Len(t, list, 0)
	})
	
	t.Run("List all items", func(t *testing.T) {
		list := service.List(0, 100)
		
		// Should contain known keys
		assert.Contains(t, list, "Hello")
		assert.Contains(t, list, "World")
		assert.Contains(t, list, "Untranslated")
		assert.Contains(t, list, "Different translation")
	})
	
	t.Run("List returns correct translations", func(t *testing.T) {
		list := service.List(0, 100)
		
		// Check specific translations
		assert.Equal(t, "Hello", list["Hello"])
		assert.Equal(t, "World", list["World"])
		// When untranslated, gotext returns the msgid itself
		assert.Equal(t, "Untranslated", list["Untranslated"])
		assert.Equal(t, "Translated value", list["Different translation"])
	})
}

func TestToOutput(t *testing.T) {
	po := createTestPo()
	service := NewPoService(*po)
	
	t.Run("Output contains headers", func(t *testing.T) {
		output := service.ToOutput()
		
		// Should contain PO file headers
		assert.Contains(t, output, "Content-Type")
		assert.Contains(t, output, "Language")
		assert.Contains(t, output, "Plural-Forms")
	})
	
	t.Run("Output contains translations", func(t *testing.T) {
		output := service.ToOutput()
		
		// Should contain msgid and msgstr entries
		assert.Contains(t, output, "msgid")
		assert.Contains(t, output, "msgstr")
		
		// Should contain specific translations
		assert.Contains(t, output, "Hello")
		assert.Contains(t, output, "World")
		assert.Contains(t, output, "Different translation")
	})
	
	t.Run("Output is valid PO format", func(t *testing.T) {
		output := service.ToOutput()
		
		// Basic PO format validation
		lines := strings.Split(output, "\n")
		require.Greater(t, len(lines), 0, "Output should have multiple lines")
		
		// Should start with msgid "" for headers
		assert.Contains(t, output, `msgid ""`)
	})
	
	t.Run("Output reflects changes", func(t *testing.T) {
		// Make a change
		service.Translate("Test Key", "Test Value")
		
		output := service.ToOutput()
		
		// Output should contain the new translation
		assert.Contains(t, output, "Test Key")
		assert.Contains(t, output, "Test Value")
	})
}

func TestIntegration(t *testing.T) {
	po := createTestPo()
	service := NewPoService(*po)
	
	t.Run("Full workflow", func(t *testing.T) {
		// 1. List untranslated items
		untranslated := service.ListAllUntranslated(0)
		initialUntranslatedCount := len(untranslated)
		assert.Greater(t, initialUntranslatedCount, 0)
		
		// 2. Translate one of them (skip the empty msgid header)
		var translatedKey string
		for key := range untranslated {
			if key != "" { // Skip the empty msgid which is the header
				translatedKey = key
				service.Translate(key, "Newly translated")
				break // Just translate the first one
			}
		}
		
		// 3. Check untranslated list is reduced and doesn't contain the translated key
		untranslatedAfter := service.ListAllUntranslated(0)
		
		// The count should be reduced (accounting for header if present)
		actualReduction := initialUntranslatedCount - len(untranslatedAfter)
		assert.Greater(t, actualReduction, 0, "Untranslated count should be reduced after translation")
		
		// The translated key should not be in the untranslated list
		if translatedKey != "" {
			assert.NotContains(t, untranslatedAfter, translatedKey)
		}
		
		// 4. List all items with pagination
		page1 := service.List(0, 2)
		page2 := service.List(2, 2)
		
		// Pages should be different
		for key := range page1 {
			assert.NotContains(t, page2, key)
		}
		
		// 5. Get output and verify it contains our translation
		output := service.ToOutput()
		assert.Contains(t, output, "Newly translated")
	})
}