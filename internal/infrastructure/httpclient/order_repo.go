// Package httpclient implements repository interfaces using HTTP calls to the backend API.
package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"bff-example/internal/domain/entity"
	"bff-example/internal/config"
)

// OrderRepo implements repository.OrderRepository via HTTP.
type OrderRepo struct {
	client  *http.Client
	baseURL string
}

// NewOrderRepo creates a new HTTP-backed order repository.
func NewOrderRepo(client *http.Client, cfg config.Config) *OrderRepo {
	return &OrderRepo{
		client:  client,
		baseURL: strings.TrimRight(cfg.BackendBaseURL, "/"),
	}
}

type processOrderRequest struct {
	OrderID  string `json:"orderId"`
	UserID   string `json:"userId"`
	Priority string `json:"priority"`
}

type processOrderResponse struct {
	OrderID     string `json:"orderId"`
	Status      string `json:"status"`
	EstimatedAt string `json:"estimatedAt"`
}

// ProcessOrder POSTs to the backend process_order endpoint.
func (r *OrderRepo) ProcessOrder(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	body, err := json.Marshal(processOrderRequest{
		OrderID:  order.OrderID,
		UserID:   order.UserID,
		Priority: order.Priority,
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.baseURL+"/api/process_order", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("backend request failed: %w", err)
	}
	defer closeResp(resp)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("backend returned %s for order %s", resp.Status, order.OrderID)
	}

	var out processOrderResponse
	if err := readJSONResponse(resp, &out); err != nil {
		return nil, fmt.Errorf("decode order response: %w", err)
	}

	return &entity.Order{
		OrderID:     out.OrderID,
		Status:      out.Status,
		EstimatedAt: out.EstimatedAt,
		UserID:      order.UserID,
		Priority:    order.Priority,
	}, nil
}
