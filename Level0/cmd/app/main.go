package main

import (
	"Level0/config"
	"Level0/internal/app"
	"flag"
	"log"
	"os"
)

const (
	DOCKER_MODE = "docker"
	LOCAL_MODE  = "local"
	Mode        = "mode"
)

func main() {
	modeFlag := flag.String(Mode, "", "mode for environment variables")
	flag.Parse()
	switch *modeFlag {
	case DOCKER_MODE:
		if err := config.DockerSystemVarsInit(); err != nil {
			log.Fatalf("Docker system vars initialization failed: %v", err)
		}
	case LOCAL_MODE:
		if err := config.LocalSystemVarsInit(); err != nil {
			log.Fatalf("Local system vars initialization failed: %v", err)
		}
	default:
		log.Fatalf("Invalid mode: %s", *modeFlag)
	}

	err := os.Setenv("MY_MODE", *modeFlag)
	if err != nil {
		log.Fatalf("Error setting system variable: %v", err)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	if err = app.RunApp(cfg); err != nil {
		log.Fatalf("Error running app: %v", err)
	}
}
