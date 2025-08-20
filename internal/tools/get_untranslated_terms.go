package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/i18n-mcp/internal/service"
	"github.com/rxtech-lab/i18n-mcp/internal/utils"
)

func NewGetUntranslatedTermsTool() (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("getUntranslatedTerms",
		mcp.WithDescription("Get untranslated terms from a PO file. After translating, you can use this tool to check if all terms are translated."),
		mcp.WithString("file_path",
			mcp.Required(),
			mcp.Description("The path to the .po file"),
		),
		mcp.WithString("limit",
			mcp.Description("Number of untranslated terms to return (default: 10)"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath, err := request.RequireString("file_path")
		if err != nil {
			return nil, fmt.Errorf("file_path parameter is required: %w", err)
		}

		// Get limit parameter, default to 10
		limitStr := request.GetString("limit", "10")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid limit value: %v", err)), nil
		}

		// Parse the PO file
		po, err := utils.ParsePoFile(filePath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error parsing PO file: %v", err)), nil
		}

		// Create PoService instance
		poService := service.NewPoService(po)

		// Get untranslated terms
		untranslatedTerms := poService.ListAllUntranslated(limit)

		// Create result object
		result := map[string]any{
			"file_path":          filePath,
			"limit":              limit,
			"count":              len(untranslatedTerms.Terms),
			"language":           untranslatedTerms.Language,
			"untranslated_terms": untranslatedTerms.Terms,
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error formatting result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}

	return tool, handler
}
