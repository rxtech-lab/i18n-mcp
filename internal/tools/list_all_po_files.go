package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/i18n-mcp/internal/utils"
)

func NewListAllPoFilesTool() (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("listAllPoFiles",
		mcp.WithDescription("List all .po files in the given directory with language information"),
		mcp.WithString("directory",
			mcp.Required(),
			mcp.Description("The directory path to scan for .po files"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		directory, err := request.RequireString("directory")
		if err != nil {
			return nil, fmt.Errorf("directory parameter is required: %w", err)
		}

		// Scan for PO files with language information
		poFilesInfo, err := utils.ScanPoFilesWithInfo(directory)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error scanning for PO files: %v", err)), nil
		}

		// Create result object
		result := map[string]interface{}{
			"directory": directory,
			"count":     len(poFilesInfo),
			"files":     poFilesInfo,
		}

		resultJSON, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error formatting result: %v", err)), nil
		}

		return mcp.NewToolResultText(string(resultJSON)), nil
	}

	return tool, handler
}
