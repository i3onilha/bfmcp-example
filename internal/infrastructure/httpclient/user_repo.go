package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"bfm-example/internal/config"
	"bfm-example/internal/domain/entity"
	"bfm-example/internal/domain/repository"
	"bfm-example/pkg/httpjson"
)

// UserRepo implements repository.UserRepository via HTTP.
type UserRepo struct {
	client  *http.Client
	baseURL string
}

// NewUserRepo creates a new HTTP-backed user repository.
func NewUserRepo(client *http.Client, cfg config.Config) *UserRepo {
	return &UserRepo{
		client:  client,
		baseURL: strings.TrimRight(cfg.BackendBaseURL, "/"),
	}
}

// GetByID GETs the backend user resource.
func (r *UserRepo) GetByID(ctx context.Context, userID string) (*entity.User, error) {
	reqURL := r.baseURL + "/api/users/" + url.PathEscape(userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("backend request failed: %w", err)
	}
	defer httpjson.DrainAndClose(resp)

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("%w: %s", repository.ErrUserNotFound, userID)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("backend returned %s for user %s", resp.Status, userID)
	}

	var out entity.User
	if err := httpjson.DecodeJSON(resp, &out); err != nil {
		return nil, fmt.Errorf("decode user response: %w", err)
	}
	return &out, nil
}
