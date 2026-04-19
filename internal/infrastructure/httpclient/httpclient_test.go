package httpclient_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"bff-example/internal/domain/entity"
	"bff-example/internal/domain/repository"
	"bff-example/internal/infrastructure/httpclient"
	"bff-example/internal/config"
)

func TestUserRepo_GetByID_success(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/users/u1" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"id":"u1","name":"Alice","email":"a@example.com"}`)
	}))
	t.Cleanup(srv.Close)

	repo := httpclient.NewUserRepo(srv.Client(), config.Config{BackendBaseURL: srv.URL})
	user, err := repo.GetByID(context.Background(), "u1")
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if user.ID != "u1" || user.Email != "a@example.com" {
		t.Fatalf("user = %#v", user)
	}
}

func TestUserRepo_GetByID_notFound(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"nope"}`, http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)

	repo := httpclient.NewUserRepo(srv.Client(), config.Config{BackendBaseURL: srv.URL})
	_, err := repo.GetByID(context.Background(), "missing")
	if err == nil || !errors.Is(err, repository.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUserRepo_GetByID_responseTooLarge(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Slightly over 1 MiB forces readJSONResponse to reject after bounded read.
		_, _ = w.Write(make([]byte, 1<<20+1024))
	}))
	t.Cleanup(srv.Close)

	repo := httpclient.NewUserRepo(srv.Client(), config.Config{BackendBaseURL: srv.URL})
	_, err := repo.GetByID(context.Background(), "u1")
	if err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("expected size error, got %v", err)
	}
}

func TestOrderRepo_ProcessOrder_success(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/process_order" || r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"orderId":"o1","status":"ok","estimatedAt":"2026-01-01T00:00:00Z"}`)
	}))
	t.Cleanup(srv.Close)

	repo := httpclient.NewOrderRepo(srv.Client(), config.Config{BackendBaseURL: srv.URL})
	out, err := repo.ProcessOrder(context.Background(), &entity.Order{
		OrderID:  "o1",
		UserID:   "u1",
		Priority: entity.PriorityNormal,
	})
	if err != nil {
		t.Fatalf("ProcessOrder: %v", err)
	}
	if out.Status != "ok" || out.OrderID != "o1" {
		t.Fatalf("order = %#v", out)
	}
}

func TestOrderRepo_ProcessOrder_errorStatus(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad", http.StatusBadRequest)
	}))
	t.Cleanup(srv.Close)

	repo := httpclient.NewOrderRepo(srv.Client(), config.Config{BackendBaseURL: srv.URL})
	_, err := repo.ProcessOrder(context.Background(), &entity.Order{OrderID: "o1", UserID: "u1"})
	if err == nil || !strings.Contains(err.Error(), "backend returned") {
		t.Fatalf("expected backend error, got %v", err)
	}
}
