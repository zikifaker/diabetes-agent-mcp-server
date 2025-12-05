package main

import (
	"context"
	"diabetes-agent-mcp-server/config"
	"diabetes-agent-mcp-server/dao"
	"diabetes-agent-mcp-server/server"
	"log/slog"
	"os"
)

func main() {
	setSysLog()

	ctx := context.Background()
	defer dao.Driver.Close(ctx)

	s := server.NewHTTPServer()
	if err := s.Start(":" + config.Cfg.Server.Port); err != nil {
		slog.Error("Failed to start MCP server", "err", err)
	}
}

func setSysLog() {
	var level slog.Leveler
	switch config.Cfg.Server.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})))
}
