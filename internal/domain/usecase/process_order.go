// Package usecase implements the ProcessOrder business logic.
package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"

	"bff-example/internal/domain/entity"
	"bff-example/internal/domain/repository"
)

// ProcessOrder handles the order processing workflow:
// 1. Validates the user exists
// 2. Processes the order
// 3. Enriches the response with user data and BFF metadata.
type ProcessOrder struct {
	userRepo  repository.UserRepository
	orderRepo repository.OrderRepository
	validator *validator.Validate
}

// NewProcessOrder creates a new ProcessOrder use case.
func NewProcessOrder(userRepo repository.UserRepository, orderRepo repository.OrderRepository, validator *validator.Validate) *ProcessOrder {
	return &ProcessOrder{
		userRepo:  userRepo,
		orderRepo: orderRepo,
		validator: validator,
	}
}

// Execute processes an order through the BFF workflow.
func (p *ProcessOrder) Execute(ctx context.Context, input ProcessOrderInput) (*ProcessOrderOutput, error) {
	if err := p.validator.Struct(&input); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidOrderInput, err)
	}

	// Step 1: Validate user exists.
	user, err := p.userRepo.GetByID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("user validation failed: %w", err)
	}

	// Step 2: Process the order.
	order := &entity.Order{
		OrderID:  input.OrderID,
		UserID:   input.UserID,
		Priority: input.Priority,
	}
	processedOrder, err := p.orderRepo.ProcessOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("order processing failed: %w", err)
	}

	// Step 3: Build enriched output.
	return &ProcessOrderOutput{
		OrderID:     processedOrder.OrderID,
		Status:      processedOrder.Status,
		EstimatedAt: processedOrder.EstimatedAt,
		UserEmail:   user.Email,
		UserName:    user.Name,
		BFF: BFFMeta{
			ProcessedBy: "bff-server",
			ProcessedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}, nil
}
