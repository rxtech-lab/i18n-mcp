package mcp

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rxtech-lab/i18n-mcp/internal/service"
	"github.com/rxtech-lab/i18n-mcp/internal/tools"
)

type MCPServer struct {
	server    *server.MCPServer
	poService *service.PoService
}

func NewMCPServer(poService *service.PoService) *MCPServer {
	mcpServer := &MCPServer{
		poService: poService,
	}
	mcpServer.InitializeTools()
	return mcpServer
}

func (s *MCPServer) InitializeTools() {
	srv := server.NewMCPServer(
		"PO Translation MCP Server",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Initialize all PO translation tools

	// 1. List all PO files tool
	listPoFilesTool, listPoFilesHandler := tools.NewListAllPoFilesTool()
	srv.AddTool(listPoFilesTool, listPoFilesHandler)

	// 2. Get untranslated terms tool
	getUntranslatedTool, getUntranslatedHandler := tools.NewGetUntranslatedTermsTool()
	srv.AddTool(getUntranslatedTool, getUntranslatedHandler)

	// 3. Look up translation tool
	lookUpTool, lookUpHandler := tools.NewLookUpTranslationTool()
	srv.AddTool(lookUpTool, lookUpHandler)

	// 4. Translate tool
	translateTool, translateHandler := tools.NewTranslateTool()
	srv.AddTool(translateTool, translateHandler)

	s.server = srv
}

func (s *MCPServer) Start() error {
	return server.ServeStdio(s.server)
}
