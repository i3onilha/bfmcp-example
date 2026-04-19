// Package repository defines interfaces for data access.
package repository

import (
	"context"

	"bfm-example/internal/domain/entity"
)

// UserRepository defines the contract for user data access.
type UserRepository interface {
	GetByID(ctx context.Context, userID string) (*entity.User, error)
}
