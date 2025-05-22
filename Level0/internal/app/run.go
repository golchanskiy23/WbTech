package app

import (
	"Level0/config"
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

const (
	SQLInitFile = "init.sql"
	Channel     = "subject"
	QueueGroup  = "queue_group"
)

func InitDB(db *postgres.DatabaseSource, path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	if _, err = db.Pool.Exec(context.Background(), string(file)); err != nil {
		return fmt.Errorf("error executing sql: %v", err)
	}
	givenOrder, err := utils.GetGivenOrder()
	if err != nil {
		return fmt.Errorf("error getting givenOrder: %v", err)
	}

	if err = postgres.AddOrdersToDB(db, givenOrder); err != nil {
		return fmt.Errorf("error adding orders to database: %v", err)
	}
	return nil
}

func CheckIsDBEmpty(db database.DatabaseRepository, ctx context.Context) bool {
	var count int
	err := db.DB.Pool.QueryRow(ctx, `
	SELECT COUNT(*) 
	FROM information_schema.tables 
	WHERE table_schema = 'public'
`).Scan(&count)

	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		return false
	}
	return true
}

func RunApp(cfg *config.Config) {
	db, err := postgres.NewStorage(postgres.GetConnection(&cfg.Database), postgres.SetMaxPoolSize(cfg.Database.MaxPoolSize))
	if err = db.Pool.Ping(context.Background()); err != nil {
		log.Fatalf("Error during creation database source: %v", err)
	}
	defer func(db *postgres.DatabaseSource) {
		db.Close()
	}(db)

	natsSrc, err := nats.NewNatsSource(&cfg.NatsStreaming)
	if err != nil {
		log.Fatalf("Nats source error: %v", err)
	}

	defer natsSrc.Close()

	pgRepository := database.CreateNewDBRepository(db)
	natsRepository := natsstreaming.CreateNewNatsStreamingRepository(natsSrc)
	// для docker
	/*if CheckIsDBEmpty(pgRepository, context.Background()) {
		if err = InitDB(db, SQLInitFile); err != nil {
			log.Fatal(err)
			return
		}
	}*/
	cacheRepository, err := cache.CreateNewCacheRepository(pgRepository)
	if err != nil {
		log.Fatalf("Error during creation of repository: %v", err)
	}

	if cacheRepository.IsEmpty() {
		if err = InitDB(db, SQLInitFile); err != nil {
			log.Fatal(err)
			return
		}
	}

	natsService := service.CreateNewNatsService(pgRepository, natsRepository, cacheRepository)
	subscription, err := natsService.StartSubscribing(Channel, QueueGroup)
	defer subscription.Close()
	if err != nil {
		log.Fatalf("Failed to start nats service: %v", err)
	}

	orderService := service.CreateNewOrderService(cacheRepository)
	orderController := controller.CreateNewOrderController(orderService)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	publisherChannel := make(chan error, 1)
	go func() {
		publisherChannel <- ExecutePublisher()
		close(publisherChannel)
	}()

	select {
	case err = <-publisherChannel:
		if err != nil {
			log.Fatalf("Failed to execute publisher: %v", err)
		}
	case <-ctx.Done():
		go func() {
			if err = <-publisherChannel; err != nil {
				log.Fatalf("Failed during get some messages: %v", err)
			}
		}()
	}

	server := server.NewServer(orderController,
		server.SetReadTimeout(6*time.Second),
		server.SetWriteTimeout(6*time.Second),
		server.SetAddr(),
		server.SetShutdownTimeout(cfg.Server.ShutdownTimeout))
	server.GracefulShutdown()
}
