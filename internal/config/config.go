package config

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	BackendBaseURL string
	Port           string
}

func Load() (Config, error) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	_ = v.BindEnv("port", "PORT")
	_ = v.BindEnv("backend.base_url", "BACKEND_BASE_URL")

	v.SetDefault("backend.base_url", "http://localhost:8082")
	v.SetDefault("port", ":8081")

	cfg := Config{
		BackendBaseURL: strings.TrimSpace(v.GetString("backend.base_url")),
		Port:           v.GetString("port"),
	}
	if cfg.BackendBaseURL == "" {
		return Config{}, errors.New("missing config: backend.base_url")
	}
	if cfg.Port == "" {
		return Config{}, errors.New("missing config: port")
	}

	return cfg, nil
}
