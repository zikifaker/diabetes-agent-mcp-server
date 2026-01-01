package server

import (
	"diabetes-agent-mcp-server/middleware"
	"diabetes-agent-mcp-server/tools"
	_ "embed"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	serverName    = "Diabetes Knowledge Search MCP Server"
	serverVersion = "1.0.0"
)

//go:embed prompts/search_diabetes_kg/query.txt
var searchDiabetesKGQueryDesc string

func NewHTTPServer() *server.StreamableHTTPServer {
	hooks := &server.Hooks{}

	// 注册工具调用成功后将调用结果推送给客户端的 hook
	hooks.AddAfterCallTool(pushCallToolResult)

	s := server.NewMCPServer(serverName, serverVersion,
		server.WithToolCapabilities(true),
		server.WithToolHandlerMiddleware(middleware.AuthMiddleware),
		server.WithHooks(hooks),
	)

	registerTools(s)

	return server.NewStreamableHTTPServer(s)
}

func registerTools(s *server.MCPServer) {
	s.AddTool(
		mcp.NewTool("search_diabetes_knowledge_graph",
			mcp.WithDescription(`
				Search professional information about diabetes guidelines, medications, diagnostics, and treatments. 
				Returns structured data from knowledge graph (entities and relationships). 
				All results are sorted by relevance score in descending order.
			`),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description(searchDiabetesKGQueryDesc),
			),
			mcp.WithNumber("limit",
				mcp.Min(10),
				mcp.Max(20),
				mcp.Description("Number of results to return"),
			),
		),
		tools.SearchDiabetesKnowledgeGraph,
	)

	s.AddTool(
		mcp.NewTool("search_user_knowledge_base",
			mcp.WithDescription(`
				Search the user's private knowledge base containing personal documents and information across various domains. 
				Use this tool when you need to find specific information from the user's personal knowledge collection, especially when general knowledge is insufficient.
			`),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description("Professionally crafted diabetes-related query for RAG system, including key medical terminology and clinical context"),
			),
			mcp.WithNumber("limit",
				mcp.Min(10),
				mcp.Max(20),
				mcp.Description("Number of results to return"),
			),
		),
		tools.SearchUserKnowledgeBase,
	)
}
