package app

import (
	"Level0/config"
	"Level0/internal/repository/cache"
	"Level0/internal/repository/database"
	"Level0/internal/repository/natsstreaming"
	"Level0/pkg/nats"
	"Level0/pkg/postgres"
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func InitDB(db *postgres.DatabaseSource, path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(1)
		log.Fatal(err)
		return err
	}
	_, err = db.Pool.Exec(context.Background(), string(file))
	if err != nil {
		fmt.Println(2)
		log.Fatal(err)
		return err
	}
	return nil
}

func RunApp(cfg *config.Config) {
	// инициализировать DB, включая туда данные из model.json(для этой задачи)
	// если нужно изменить колонки таблицы - миграции, а не создание таблицы с нуля
	fmt.Printf("Config for database: %v\n", cfg.Database)
	db, err := postgres.NewStorage(postgres.GetConnection(&cfg.Database), postgres.SetMaxPoolSize(cfg.Database.MaxPoolSize))
	err = db.Pool.Ping(context.Background())
	if err != nil {
		//поменять ошибку
		log.Fatal(err)
	}
	fmt.Printf("Connected to database: %v", db)
	defer func(db *postgres.DatabaseSource) {
		db.Close()
	}(db)

	err = InitDB(db, "init.sql")
	if err != nil {
		log.Fatal(err)
		return
	}

	natsSrc, err := nats.NewNatsSource(&cfg.NatsStreaming)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Connected to nats source: %v", natsSrc)
	defer natsSrc.Close()

	pgRepository := database.CreateNewDBRepository(db)
	natsRepository := natsstreaming.CreateNewNatsStreamingRepository(natsSrc)
	fmt.Printf("DataSources created: %v %v\n", natsRepository, pgRepository)

	cacheRepository, err := cache.CreateNewCacheRepository(pgRepository)
	if err != nil {
		fmt.Printf("Error during creation of repository: %v", err)
		return
	}
	fmt.Println(cacheRepository)

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
