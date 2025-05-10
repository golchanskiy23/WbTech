package app

import (
	"Level0/config"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func RunApp(cfg *config.Config) {
	// используем элементы из кофигурации
	// driver_name, name, password, host, port, db, ssl_mode
	dsn := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s", "postgres",
		os.Getenv("POSTRGES_UNSAFE_USERNAME"),
		os.Getenv("POSTRGES_UNSAFE_PASSWORD"),
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)
	// некастомная установка соединения
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("1")
		log.Fatalf("Failed to open DB: %v", err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("2")
		fmt.Println(err)
	}

	// инициализируем postgres (каким способом - ?)(возможно потребуется более тонкая настройка)
	// аналогичная инициализация для nats
	nats, err := stan.Connect(os.Getenv("CLUSTER_ID"), os.Getenv("SECOND_CLIENT_ID"), stan.NatsURL(cfg.NatsStreaming.URL))
	if err != nil {
		fmt.Println("3")
		fmt.Println(err)
		return
	}

	//  хранилища данных
	pgStorage := CreatNewPgStorage(db)
	cacheStorage := CreateNewCacheStorage(pgStorage)
	natsStorage := CreateNewNatsStorage(nats)

	// использует все хранилища( скорее всего будет какая-то инициализация)
	natsService := CreateNewService(pgStorage, cacheStorage, natsStorage)
	// подписка на канал
	// channel, queue_group(для единовременной отправки или получения сообщений)
	err = natsService.StartNatsService()
	if err != nil {
		fmt.Println("4")
		fmt.Printf("Failed to start nats service: %v", err)
	}

	// бизнес-логика обработки заказов
	orderService := CreateNewOrderService(pgStorage, cacheStorage)
	// принимает и роутит запросы по нужным ручкам
	orderController := CreateNewOrderController(orderService)
	// и для сервера(+кастомизация с нужными параметрами)
	server := http.Server{}
	// запуск сервера , хэндлер которого - контроллер, роутящий запросы
	http.ListenAndServe(fmt.Sprintf("localhost:%s", cfg.Server.Addr), orderController)

	osInterruptChan := make(chan os.Signal, 1)
	signal.Notify(osInterruptChan, os.Interrupt)
	select {
	case <-osInterruptChan:
		log.Fatal("OS interruption")
	case <-server.Done():
		log.Fatal("Server threw an error")
	}

	// context - (?)
	err = server.Shutdown(context.Background())
	if err != nil {
		log.Fatal("Server shutdown failure")
	}

	log.Println("Server shutdown complete")
}
