package config

import (
	"fmt"
	"github.com/joho/godotenv"
)

func DockerSystemVarsInit() error {
	if err := godotenv.Load(".env.container"); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}
	return nil
}

func LocalSystemVarsInit() error {
	if err := godotenv.Load(".env.local"); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}
	return nil
}
