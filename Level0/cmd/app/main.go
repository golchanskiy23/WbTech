package main

import (
	"Level0/config"
	"Level0/internal/app"
)

// initConfig from config-package
// config for db
// config for nats streaming
// config for http-server
// config for app - ?
func init() {

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
