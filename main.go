// Copyright 2025 The Go MCP SDK Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package main is a demo client that connects to the BFF server and calls its tools.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"bff-example/pkg/headerforward"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// === Infrastructure Layer: Backend HTTP API ===
	go runBackendHTTPServer()
	time.Sleep(100 * time.Millisecond)

	ctx := context.Background()
	ctx = context.WithValue(ctx, headerforward.ContextKey{}, headerforward.FilterHeaders(http.Header{
		"X-Tenant-Id":      []string{"tenant-123"},
		"X-Correlation-Id": []string{"corr-abc"},
	}))
	// Start the backend and BFF servers.
	backendAddr := "localhost:8081"
	// Connect to the BFF server.
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "demo-client",
		Version: "1.0.0",
	}, nil)
	clientSession, err := client.Connect(ctx, &mcp.StreamableClientTransport{
		Endpoint:   fmt.Sprintf("http://%s/mcp", backendAddr),
		HTTPClient: headerforward.NewClient(),
	}, nil)
	if err != nil {
		log.Fatalf("Failed to connect to BFF: %v", err)
	}
	defer clientSession.Close()

	// Call the BFF process_order tool with typed parameters.
	orderResult, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name: "process_order",
		Arguments: map[string]any{
			"orderId":  "ORD-42",
			"userId":   "u1",
			"priority": "high",
		},
	})
	if err != nil {
		log.Fatalf("Failed to process order: %v", err)
	}
	buf, err := json.Marshal(orderResult)
	fmt.Println(string(buf))
}

// runBackendHTTPServer starts a mock backend REST API for testing.
// In production, this would be a separate service.
func runBackendHTTPServer() {
	mux := http.NewServeMux()

	type processOrderBody struct {
		OrderID  string `json:"orderId"`
		UserID   string `json:"userId"`
		Priority string `json:"priority"`
	}
	type orderResult struct {
		OrderID     string `json:"orderId"`
		Status      string `json:"status"`
		EstimatedAt string `json:"estimatedAt"`
	}

	mux.HandleFunc("POST /api/process_order", func(w http.ResponseWriter, r *http.Request) {
		var body processOrderBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
			return
		}
		status := "confirmed"
		if body.Priority == "high" {
			status = "expedited"
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(orderResult{
			OrderID:     body.OrderID,
			Status:      status,
			EstimatedAt: "2026-04-18T10:00:00Z",
		})
	})

	type userInfo struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	mux.HandleFunc("GET /api/users/{userID}", func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("userID")
		users := map[string]userInfo{
			"u1": {"u1", "Alice", "alice@example.com"},
			"u2": {"u2", "Bob", "bob@example.com"},
		}
		user, ok := users[userID]
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"user not found"}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(user)
	})

	log.Println("Starting Backend HTTP Server on :8082")
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatal(err)
	}
}
