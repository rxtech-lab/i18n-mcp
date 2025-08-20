package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/i18n-mcp/internal/service"
	"github.com/rxtech-lab/i18n-mcp/internal/utils"
)

func NewTranslateTool() (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("translate",
		mcp.WithDescription("Translate terms in a PO file and save the changes. You can translate multiple terms at once or updating the existing translation."),
		mcp.WithString("file_path",
			mcp.Required(),
			mcp.Description("The path to the .po file"),
		),
		mcp.WithString("translations",
			mcp.Required(),
			mcp.Description("JSON object with translations where keys are term keys and values are translations"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath, err := request.RequireString("file_path")
		if err != nil {
			return nil, fmt.Errorf("file_path parameter is required: %w", err)
		}

		translationsStr, err := request.RequireString("translations")
		if err != nil {
			return nil, fmt.Errorf("translations parameter is required: %w", err)
		}

		// Parse translations JSON
		var translations map[string]string
		if err := json.Unmarshal([]byte(translationsStr), &translations); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid translations JSON: %v", err)), nil
		}

		// Parse the PO file
		po, err := utils.ParsePoFile(filePath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error parsing PO file: %v", err)), nil
		}

		// Create PoService instance
		poService := service.NewPoService(po)

		// Apply translations
		translatedCount := 0
		for key, value := range translations {
			poService.Translate(key, value)
			translatedCount++
		}

		// Get the updated PO file content
		updatedContent := poService.ToOutput()

		// Write the updated content back to the file
		err = os.WriteFile(filePath, []byte(updatedContent), 0644)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error writing to PO file: %v", err)), nil
		}

		// Create result object
		result := map[string]interface{}{
			"file_path":        filePath,
			"translated_count": translatedCount,
			"translations":     translations,
			"message":          fmt.Sprintf("Successfully translated %d terms and saved to %s", translatedCount, filePath),
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error formatting result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}

	return tool, handler
}
