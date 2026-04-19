package config_test

import (
	"testing"

	"bfm-example/internal/config"
)

func TestLoad_envOverrides(t *testing.T) {
	t.Setenv("PORT", ":5999")
	t.Setenv("BACKEND_BASE_URL", "https://backend.example/")
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != ":5999" {
		t.Fatalf("Port = %q want :5999", cfg.Port)
	}
	if cfg.BackendBaseURL != "https://backend.example/" {
		t.Fatalf("BackendBaseURL = %q", cfg.BackendBaseURL)
	}
}
