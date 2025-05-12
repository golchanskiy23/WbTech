package app

import (
	"Level0/config"
	"Level0/internal/repository/database"
	"Level0/internal/repository/natsstreaming"
	"Level0/pkg/nats"
	"Level0/pkg/postgres"
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

func RunApp(cfg *config.Config) {
	fmt.Printf("Config for database: %v\n", cfg.Database)
	db, err := postgres.NewStorage(postgres.GetConnection(&cfg.Database), postgres.SetMaxPoolSize(cfg.Database.MaxPoolSize))
	err = db.Pool.Ping(context.Background())
	if err != nil {
		//поменять ошибку
		log.Fatal(err)
	}
	fmt.Printf("Connected to database: %v", db)
	defer db.Close()

	natsSrc, err := nats.NewNatsSource(&cfg.NatsStreaming)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Connected to nats source: %v", natsSrc)
	defer natsSrc.Close()

	pgRepository := database.CreateNewDBRepository(db)
	natsRepository := natsstreaming.CreateNewNatsStreamingRepository(natsSrc)
	fmt.Printf("DataSources created: %v %v\n", natsRepository, pgRepository)
	/*cacheStorage := CreateNewCacheStorage(pgStorage)
	 */

	// использует все хранилища( скорее всего будет какая-то инициализация)
	/*natsService := CreateNewService(pgStorage, cacheStorage, natsStorage)
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
	err = server.Shutdown_()
	if err != nil {
		log.Fatal("Server shutdown failure")
	}

	log.Println("Server shutdown complete")
	*/
}
