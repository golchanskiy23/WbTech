package main

import (
	"Level0/config"
	"Level0/internal/app"
)

// Подумать над миграциями
// initConfig from config-package
// config for db
// config for nats streaming
// config for http-server
// config for app - ?

func main() {
	config.SystemVarsInit()
	cfg, err := config.NewConfig()
	if err != nil {
		// какой-то обработчик ошибок
		return
	}
	app.RunApp(cfg)
}
