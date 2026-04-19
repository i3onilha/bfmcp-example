// Package repository defines interfaces for data access.
package repository

import (
	"context"

	"bfm-example/internal/domain/entity"
)

// OrderRepository defines the contract for order data access.
type OrderRepository interface {
	ProcessOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
}
