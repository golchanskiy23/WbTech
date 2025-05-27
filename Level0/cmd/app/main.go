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
)

func main() {
	modeFlag := flag.String("mode", "", "mode for environment variables")
	flag.Parse()
	if *modeFlag == DOCKER_MODE {
		if err := config.DockerSystemVarsInit(); err != nil {
			log.Fatalf("Docker system vars initialization failed: %v", err)
		}
	} else if *modeFlag == LOCAL_MODE {
		if err := config.LocalSystemVarsInit(); err != nil {
			log.Fatalf("Local system vars initialization failed: %v", err)
		}
	} else {
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
	app.RunApp(cfg)
}
