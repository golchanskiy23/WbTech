package config

import (
	"fmt"
	"github.com/joho/godotenv"
)

func SystemVarsInit() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}
	return nil
}
