package usecase_test

import (
	"context"
	"errors"
	"testing"

	"bff-example/internal/domain/entity"
	"bff-example/internal/domain/repository"
	"bff-example/internal/domain/usecase"
)

type fakeUserRepo struct {
	user *entity.User
	err  error
}

func (f *fakeUserRepo) GetByID(ctx context.Context, userID string) (*entity.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.user, nil
}

type fakeOrderRepo struct {
	order *entity.Order
	err   error
}

func (f *fakeOrderRepo) ProcessOrder(ctx context.Context, order *entity.Order) (*entity.Order, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.order, nil
}

func TestProcessOrder_Execute_success(t *testing.T) {
	t.Parallel()
	uc := usecase.NewProcessOrder(
		&fakeUserRepo{user: &entity.User{ID: "u1", Name: "A", Email: "a@x"}},
		&fakeOrderRepo{order: &entity.Order{OrderID: "o1", Status: "ok", EstimatedAt: "t"}},
	)
	out, err := uc.Execute(context.Background(), usecase.ProcessOrderInput{
		OrderID:  "o1",
		UserID:   "u1",
		Priority: entity.PriorityHigh,
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if out.UserEmail != "a@x" || out.Status != "ok" {
		t.Fatalf("output = %#v", out)
	}
}

func TestProcessOrder_Execute_validation(t *testing.T) {
	t.Parallel()
	uc := usecase.NewProcessOrder(&fakeUserRepo{}, &fakeOrderRepo{})
	_, err := uc.Execute(context.Background(), usecase.ProcessOrderInput{
		OrderID:  "",
		UserID:   "u1",
		Priority: entity.PriorityHigh,
	})
	if !errors.Is(err, usecase.ErrInvalidOrderInput) {
		t.Fatalf("want ErrInvalidOrderInput, got %v", err)
	}
}

func TestProcessOrder_Execute_userNotFound(t *testing.T) {
	t.Parallel()
	uc := usecase.NewProcessOrder(
		&fakeUserRepo{err: repository.ErrUserNotFound},
		&fakeOrderRepo{},
	)
	_, err := uc.Execute(context.Background(), usecase.ProcessOrderInput{
		OrderID:  "o1",
		UserID:   "u1",
		Priority: entity.PriorityNormal,
	})
	if !errors.Is(err, repository.ErrUserNotFound) {
		t.Fatalf("errors.Is user not found: %v", err)
	}
}
