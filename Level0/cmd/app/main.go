package main

import (
	"Level0/config"
	"Level0/internal/app"
	"github.com/joho/godotenv"
	"log"
)

// Подумать над миграциями
// initConfig from config-package
// config for db
// config for nats streaming
// config for http-server
// config for app - ?
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		// какой-то обработчик ошибок
		return
	}
	// запуск самой программы: предварительная инициализация,
	// запуск прослушки сообщений через каналы(get data , interrupting server)
	// shutdown server
	app.RunApp(cfg)
}
