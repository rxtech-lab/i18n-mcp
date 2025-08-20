package tools

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
)

// getTextContent extracts text content from a CallToolResult
func getTextContent(t *testing.T, result *mcp.CallToolResult) string {
	require.NotNil(t, result)
	require.NotEmpty(t, result.Content)

	// Use the AsTextContent helper function to convert the Content interface
	textContent, ok := mcp.AsTextContent(result.Content[0])
	require.True(t, ok, "Expected TextContent")

	return textContent.Text
}

// makeRequest creates a CallToolRequest with the given arguments
func makeRequest(args map[string]interface{}) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: args,
		},
	}
}
