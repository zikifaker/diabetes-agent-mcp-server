package server

import (
	"diabetes-agent-mcp-server/tool"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	serverName    = "Diabetes Agent MCP Server"
	serverVersion = "1.0.0"
)

func NewHTTPServer() *server.StreamableHTTPServer {
	s := server.NewMCPServer(serverName, serverVersion,
		server.WithToolCapabilities(true),
	)

	s.AddTool(
		mcp.NewTool("diabetes_knowledge_base_search",
			mcp.WithDescription("search professional information about diabetes guidelines, medications, diagnostics, and treatments"),
			mcp.WithString("query", mcp.Required(), mcp.Description("search query")),
			mcp.WithNumber("limit", mcp.DefaultNumber(tool.DefaultSearchResultLimit), mcp.Description("search result limit, you'd better use an even number")),
		),
		tool.SearchDiabetesKnowledgeBase,
	)

	return server.NewStreamableHTTPServer(s)
}
