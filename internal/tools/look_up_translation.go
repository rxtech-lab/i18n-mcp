package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/i18n-mcp/internal/service"
	"github.com/rxtech-lab/i18n-mcp/internal/utils"
)

func NewLookUpTranslationTool() (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("lookUpTranslation",
		mcp.WithDescription("Search for a term key and return the translated value from a PO file. Use this tool to look up the previous translation of a term."),
		mcp.WithString("file_path",
			mcp.Required(),
			mcp.Description("The path to the .po file"),
		),
		mcp.WithString("search_term",
			mcp.Required(),
			mcp.Description("The term key to search for"),
		),
		mcp.WithString("page_size",
			mcp.Description("Number of results to return per page (default: 10)"),
		),
		mcp.WithString("page",
			mcp.Description("Page number for pagination (default: 1)"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath, err := request.RequireString("file_path")
		if err != nil {
			return nil, fmt.Errorf("file_path parameter is required: %w", err)
		}

		searchTerm, err := request.RequireString("search_term")
		if err != nil {
			return nil, fmt.Errorf("search_term parameter is required: %w", err)
		}

		// Get pagination parameters
		pageSizeStr := request.GetString("page_size", "10")
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid page_size value: %v", err)), nil
		}

		pageStr := request.GetString("page", "1")
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid page value: %v", err)), nil
		}

		// Parse the PO file
		po, err := utils.ParsePoFile(filePath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error parsing PO file: %v", err)), nil
		}

		// Create PoService instance
		poService := service.NewPoService(po)

		// Get all translations for searching
		// Using a large number to get all translations for filtering
		allTranslations := poService.List(0, 100000)

		// Filter translations that contain the search term
		matchingTranslations := make(map[string]string)
		for msgid, msgstr := range allTranslations {
			if strings.Contains(strings.ToLower(msgid), strings.ToLower(searchTerm)) {
				matchingTranslations[msgid] = msgstr
			}
		}

		// Apply pagination to matching results
		skip := (page - 1) * pageSize
		paginatedResults := make(map[string]string)
		count := 0
		taken := 0

		for msgid, msgstr := range matchingTranslations {
			if count < skip {
				count++
				continue
			}
			if taken >= pageSize {
				break
			}
			paginatedResults[msgid] = msgstr
			taken++
			count++
		}

		// Create result object
		result := map[string]interface{}{
			"file_path":     filePath,
			"search_term":   searchTerm,
			"page":          page,
			"page_size":     pageSize,
			"total_matches": len(matchingTranslations),
			"translations":  paginatedResults,
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error formatting result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}

	return tool, handler
}
