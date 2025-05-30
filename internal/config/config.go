package config

import (
	"github.com/joho/godotenv"
)

func Load(path string) error {
	// Try to load .env file but don't fail if it doesn't exist

	return godotenv.Load(path)
}
