// Package transport provides MCP server setup and configuration.
package transport

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"bfm-example/internal/adapter/handler"
	"bfm-example/internal/adapter/middleware"
)

// RegisterTools creates and configures an MCP BFF server with all tools.
func RegisterTools(
	orderHandler *handler.ProcessOrderHandler,
) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "bff-server",
		Version: "1.0.0",
	}, nil)

	// Add logging middleware to all incoming requests.
	server.AddReceivingMiddleware(middleware.Logging())

	mcp.AddTool(server, &mcp.Tool{
		Name:        "process_order",
		Description: "Process an order through the BFF (validates user, forwards to backend, enriches response)",
	}, orderHandler.Handle)

	return server
}
