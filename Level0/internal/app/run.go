package app

import (
	"Level0/config"
	"Level0/internal/controller"
	"Level0/internal/repository/cache"
	"Level0/internal/repository/database"
	"Level0/internal/repository/natsstreaming"
	"Level0/internal/service"
	"Level0/internal/utils"
	internal_nats "Level0/pkg/nats"
	"Level0/pkg/postgres"
	"Level0/pkg/server"
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"time"
)

const (
	SQLInitFile = "init.sql"
	Channel     = "orders.getter"
	QueueGroup  = "orders_group"
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

type JetsResponse struct {
	Conn *nats.Conn
	Js   nats.JetStreamContext
	Err  error
}

func InitJetStream(cfg *config.Config) JetsResponse {
	nc, err := nats.Connect(getConnection(os.Getenv("NATS_HOST"), cfg.NatsStreaming.Port))
	if err != nil {
		return JetsResponse{Conn: nil, Err: fmt.Errorf("can't connect to JetStream: %v", err), Js: nil}
	}
	js, err := nc.JetStream()
	if err != nil {
		return JetsResponse{Conn: nil, Err: fmt.Errorf("can't get JetStream context: %v", err), Js: nil}
	}

	_, err = js.StreamInfo("ORDER_STREAM")
	if err == nats.ErrStreamNotFound {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     "ORDER_STREAM",
			Subjects: []string{"orders.*"},
		})
		if err != nil {
			return JetsResponse{Conn: nil, Err: fmt.Errorf("can't add stream: %w", err), Js: nil}
		}
	} else if err != nil {
		return JetsResponse{Conn: nil, Err: fmt.Errorf("can't check stream: %w", err), Js: nil}
	}
	return JetsResponse{Conn: nc, Js: js, Err: nil}
}

func RunApp(cfg *config.Config) {
	db, err := postgres.NewStorage(postgres.GetConnection(&cfg.Database), postgres.SetMaxPoolSize(cfg.Database.MaxPoolSize))
	fmt.Println(postgres.GetConnection(&cfg.Database))
	if err = db.Pool.Ping(context.Background()); err != nil {
		log.Fatalf("Error during creation database source: %v", err)
	}
	defer func(db *postgres.DatabaseSource) {
		db.Close()
	}(db)

	response := InitJetStream(cfg)
	if response.Err != nil {
		log.Fatal(response.Err)
	}
	fmt.Println(response)
	natsSrc, err := internal_nats.NewNatsSource(response.Conn, response.Js)
	if err != nil {
		log.Fatalf("Nats source error: %v", err)
	}

	defer natsSrc.Close()

	pgRepository := database.CreateNewDBRepository(db)
	natsRepository := natsstreaming.NewNatsJetStreamRepository(natsSrc)
	// для docker
	orders, err := pgRepository.GetAllOrders(context.Background())
	if err != nil {
		log.Fatalf("Error during cache initialization: %v", err)
	} else if len(orders) == 0 {
		if err = InitDB(db, SQLInitFile); err != nil {
			log.Fatal(err)
			return
		}
	}

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

	natsService := service.CreateNewNatsService(pgRepository, *natsRepository, cacheRepository)
	subscription, err := natsService.StartSubscribing(Channel, QueueGroup)
	defer subscription.Unsubscribe()
	if err != nil {
		log.Fatalf("Failed to start nats service: %v", err)
	}

	orderService := service.CreateNewOrderService(cacheRepository)
	orderController := controller.CreateNewOrderController(orderService)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	publisherChannel := make(chan error, 1)
	go func(j nats.JetStreamContext, c *nats.Conn) {
		publisherChannel <- ExecutePublisher(j, c)
		close(publisherChannel)
	}(response.Js, response.Conn)

	select {
	case err = <-publisherChannel:
		if err != nil {
			log.Fatalf("Failed to execute publisher: %v", err)
		}
	case <-ctx.Done():
		go func() {
			if err = <-publisherChannel; err != nil {
				log.Fatalf("End of service work: %v", err)
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
