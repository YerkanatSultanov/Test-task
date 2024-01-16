package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

func LoadConfig() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	envFile := filepath.Join(dir, "config", ".env")
	err = godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	return nil
}
