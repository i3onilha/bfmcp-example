package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"bff-example/internal/adapter/handler"
	"bff-example/internal/domain/entity"
	"bff-example/internal/domain/usecase"
)

type panicUserRepo struct{}

func (panicUserRepo) GetByID(context.Context, string) (*entity.User, error) {
	panic("GetByID must not be called when validation fails")
}

type panicOrderRepo struct{}

func (panicOrderRepo) ProcessOrder(context.Context, *entity.Order) (*entity.Order, error) {
	panic("ProcessOrder must not be called when validation fails")
}

func TestProcessOrderHandler_validationBeforeRepos(t *testing.T) {
	t.Parallel()
	uc := usecase.NewProcessOrder(panicUserRepo{}, panicOrderRepo{})
	h := handler.NewProcessOrderHandler(uc)

	req := &mcp.CallToolRequest{Extra: &mcp.RequestExtra{}}
	_, _, err := h.Handle(context.Background(), req, &handler.ProcessOrderArgs{
		OrderID:  "",
		UserID:   "u1",
		Priority: entity.PriorityHigh,
	})
	if !errors.Is(err, usecase.ErrInvalidOrderInput) {
		t.Fatalf("want ErrInvalidOrderInput, got %v", err)
	}
}
