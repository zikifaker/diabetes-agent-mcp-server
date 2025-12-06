package tool

import (
	"context"
	"diabetes-agent-mcp-server/dao"
	"fmt"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	DefaultSearchResultLimit = 20
	Neo4jFulltextName        = "fulltextSearch"
)

// SearchDiabetesKnowledgeBase 检索糖尿病知识检索
// 先分别进行图检索和向量检索，各自召回 limit / 2 条，再将结果重排序，取 Top limit / 2 条作为最终结果
func SearchDiabetesKnowledgeBase(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := req.GetString("query", "")
	if query == "" {
		content := mcp.TextContent{
			Text: "query param is required",
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{content},
			IsError: true,
		}, nil
	}

	limit := req.GetInt("limit", DefaultSearchResultLimit)
	if limit <= 0 {
		limit = DefaultSearchResultLimit
	}

	knowledgeGraphResults, err := searchKnowledgeGraph(ctx, query, limit)
	if err != nil {
		slog.Error("Failed to search knowledge graph", "err", err)
	}

	vectorDBResults, err := searchVectorDB(ctx, query, limit)
	if err != nil {
		slog.Error("Failed to search vector DB", "err", err)
	}

	finalResults := map[string]any{
		"knowledge_graph_results": knowledgeGraphResults,
		"vector_db_results":       vectorDBResults,
	}

	slog.Debug("search diabetes knowledge base finished", "final_results", finalResults)

	return mcp.NewToolResultJSON(finalResults)
}

// 检索图数据库（DiaKG数据集）
func searchKnowledgeGraph(ctx context.Context, query string, limit int) ([]map[string]any, error) {
	session := dao.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	// 查询至少有一个关系的节点
	cypherQuery := `
	CALL db.index.fulltext.queryNodes($indexName, $query) 
	YIELD node, score
	WHERE 'Entity' IN labels(node)
	ORDER BY score DESC
	LIMIT $limit
	WITH node, score, [(node)-[r]-(related:Entity) | {
		type: type(r),
		related: related {.name, .type}
	}] AS relationships
	WHERE size(relationships) > 0
	RETURN 
		node {.name, .type} AS node,
		relationships,
		score
	`

	result, err := session.Run(ctx, cypherQuery, map[string]any{
		"indexName": Neo4jFulltextName,
		"query":     query,
		"limit":     limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute fulltext query: %v", err)
	}

	var results []map[string]any
	for result.Next(ctx) {
		record := result.Record()
		results = append(results, record.AsMap())
	}

	if err = result.Err(); err != nil {
		return nil, fmt.Errorf("failed to process search results: %v", err)
	}

	return results, nil
}

// 检索向量数据库（用户上传的知识文件）
func searchVectorDB(ctx context.Context, query string, limit int) ([]map[string]any, error) {
	return nil, nil
}
