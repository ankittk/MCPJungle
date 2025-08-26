package mcp

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
	"gorm.io/gorm"

	"github.com/mcpjungle/mcpjungle/internal/metrics"
)

// MCPService coordinates operations amongst the registry database, mcp proxy server and upstream MCP servers.
// It is responsible for maintaining data consistency and providing a unified interface for MCP operations.
type MCPService struct {
	db             *gorm.DB
	mcpProxyServer *server.MCPServer
	metrics        *metrics.MCPMetrics
}

// NewMCPService creates a new instance of MCPService.
// It initializes the MCP proxy server by loading all registered tools from the database.
func NewMCPService(db *gorm.DB, mcpProxyServer *server.MCPServer, metrics *metrics.MCPMetrics) (*MCPService, error) {
	s := &MCPService{
		db:             db,
		mcpProxyServer: mcpProxyServer,
		metrics:        metrics,
	}
	if err := s.initMCPProxyServer(); err != nil {
		return nil, fmt.Errorf("failed to initialize MCP proxy server: %w", err)
	}
	return s, nil
}
