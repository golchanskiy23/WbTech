package main

import (
	"Level0/config"
	"Level0/internal/app"
	"log"
)

func main() {
	if err := config.SystemVarsInit(); err != nil {
		log.Fatalf("System vars initialization failed: %v", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	app.RunApp(cfg)
}
