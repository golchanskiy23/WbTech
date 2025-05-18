package app

import (
	"Level0/config"
	"Level0/internal/addingutils"
	"Level0/internal/controller"
	"Level0/internal/repository/cache"
	"Level0/internal/repository/database"
	"Level0/internal/repository/natsstreaming"
	"Level0/internal/service"
	"Level0/internal/utils"
	"Level0/pkg/nats"
	"Level0/pkg/postgres"
	"Level0/pkg/server"
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
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

	err = addingutils.AddOrdersToDB(db, utils.GetGivenOrder())
	if err != nil {
		fmt.Println(3)
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

	/*err = InitDB(db, "init.sql")
	if err != nil {
		log.Fatal(err)
		return
	}*/

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

	natsService := service.CreateNewNatsService(pgRepository, natsRepository, cacheRepository)
	subscription, err := natsService.StartSubscribing("subject", "queue_group")
	defer subscription.Close()
	if err != nil {
		fmt.Println("4")
		fmt.Printf("Failed to start nats service: %v", err)
	}

	orderService := service.CreateNewOrderService(cacheRepository)
	orderController := controller.CreateNewOrderController(orderService)
	go ExecutePublisher()
	time.Sleep(2 * time.Second)

	server := server.NewServer(orderController,
		server.SetReadTimeout(6*time.Second),
		server.SetWriteTimeout(6*time.Second),
		server.SetAddr(cfg.Server.Addr),
		server.SetShutdownTimeout(cfg.Server.ShutdownTimeout))
	server.GracefulShutdown()

}
