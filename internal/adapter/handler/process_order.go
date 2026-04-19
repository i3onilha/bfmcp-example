// Package handler provides MCP tool handlers that bridge the MCP protocol
// to domain use cases.
package handler

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"bfm-example/internal/domain/usecase"
	"bfm-example/pkg/headerforward"
)

// ProcessOrderHandler handles the process_order BFF tool.
type ProcessOrderHandler struct {
	uc *usecase.ProcessOrder
}

// NewProcessOrderHandler creates a new process order handler.
func NewProcessOrderHandler(uc *usecase.ProcessOrder) *ProcessOrderHandler {
	return &ProcessOrderHandler{uc: uc}
}

// ProcessOrderArgs defines the input parameters for the process_order tool.
type ProcessOrderArgs struct {
	OrderID  string `json:"orderId" jsonschema:"The order ID to process"`
	UserID   string `json:"userId" jsonschema:"The user making the request"`
	Priority string `json:"priority" jsonschema:"Processing priority: high, normal, or low"`
}

// Handle executes the process_order tool, calling the use case and returning
// the enriched BFF response.
func (h *ProcessOrderHandler) Handle(ctx context.Context, req *mcp.CallToolRequest, args *ProcessOrderArgs) (*mcp.CallToolResult, *usecase.ProcessOrderOutput, error) {
	// Propagate allowlisted incoming HTTP headers to downstream backend calls.
	if forwarded := headerforward.FilterHeaders(req.Extra.Header); len(forwarded) > 0 {
		ctx = context.WithValue(ctx, headerforward.ContextKey{}, forwarded)
	}

	// Execute the use case.
	output, err := h.uc.Execute(ctx, usecase.ProcessOrderInput{
		OrderID:  args.OrderID,
		UserID:   args.UserID,
		Priority: args.Priority,
	})
	if err != nil {
		return nil, nil, err
	}

	return nil, output, nil
}
