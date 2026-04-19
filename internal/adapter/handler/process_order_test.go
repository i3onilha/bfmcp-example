package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"bff-example/internal/adapter/handler"
	"bff-example/internal/domain/entity"
	"bff-example/internal/domain/usecase"
	pkgvalidate "bff-example/pkg/validate"
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
	v, err := pkgvalidate.New()
	if err != nil {
		t.Fatal(err)
	}
	uc := usecase.NewProcessOrder(panicUserRepo{}, panicOrderRepo{}, v)
	h := handler.NewProcessOrderHandler(uc)

	req := &mcp.CallToolRequest{Extra: &mcp.RequestExtra{}}
	_, _, handleErr := h.Handle(context.Background(), req, &handler.ProcessOrderArgs{
		OrderID:  "",
		UserID:   "u1",
		Priority: entity.PriorityHigh,
	})
	if !errors.Is(handleErr, usecase.ErrInvalidOrderInput) {
		t.Fatalf("want ErrInvalidOrderInput, got %v", handleErr)
	}
}
