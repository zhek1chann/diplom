package config

import (
	"fmt"
	"os"
)

const (
	nctBaseURLEnv = "NCT_BASE_URL"
)

type NCTConfig interface {
	BaseURL() string
}

type nctConfig struct {
	baseURL string
}

func NewNCTConfig() (NCTConfig, error) {
	baseURL := os.Getenv(nctBaseURLEnv)
	if baseURL == "" {
		return nil, fmt.Errorf("NCT base URL is not set")
	}

	return &nctConfig{
		baseURL: baseURL,
	}, nil
}

func (c *nctConfig) BaseURL() string {
	return c.baseURL
}
