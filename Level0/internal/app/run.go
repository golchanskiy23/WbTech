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
	"github.com/go-chi/chi/v5"
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
	Stream      = "ORDER_STREAM"
)

type JetsResponse struct {
	Conn *nats.Conn
	Js   nats.JetStreamContext
	Err  error
}

type CacheResponse struct {
	Repo    *cache.CacheRepository
	Err     error
	Checker utils.Checker
}

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

func InitJetStream(cfg *config.Config) JetsResponse {
	nc, err := nats.Connect(getConnection(os.Getenv("NATS_HOST"), cfg.NatsStreaming.Port))
	if err != nil {
		return JetsResponse{Conn: nil, Err: fmt.Errorf("can't connect to JetStream: %v", err), Js: nil}
	}
	js, err := nc.JetStream()
	if err != nil {
		return JetsResponse{Conn: nil, Err: fmt.Errorf("can't get JetStream context: %v", err), Js: nil}
	}

	_, err = js.StreamInfo(Stream)
	if err == nats.ErrStreamNotFound {
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     Stream,
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

func initPostgres(cfg *config.Config) (*postgres.DatabaseSource, error) {
	db, err := postgres.NewStorage(
		postgres.GetConnection(&cfg.Database),
		postgres.SetMaxPoolSize(cfg.Database.MaxPoolSize),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init postgres: %w", err)
	}
	if err = db.Pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}
	return db, nil
}

func initDBRepository(db *postgres.DatabaseSource) database.DatabaseRepository {
	pgRepository := database.CreateNewDBRepository(db)
	return pgRepository
}

func initCache(pgRepository database.DatabaseRepository, c utils.Checker) CacheResponse {
	orders, err := pgRepository.GetAllOrders(context.Background())
	if err != nil {
		return CacheResponse{nil, fmt.Errorf("failed to get all orders: %w", err), c}
	}
	if len(orders) == 0 {
		if err = InitDB(pgRepository.DB, SQLInitFile); err != nil {
			return CacheResponse{nil, err, c}
		}
	}
	cacheRepository, err := cache.CreateNewCacheRepository(orders)
	if err != nil {
		return CacheResponse{nil, err, c}
	}
	if cacheRepository.IsEmpty() {
		if err = InitDB(pgRepository.DB, SQLInitFile); err != nil {
			return CacheResponse{nil, err, c}
		}
	}
	for _, order := range orders {
		if _, ok := c[order.OrderUID]; !ok {
			c[order.OrderUID] = struct{}{}
		}
	}
	return CacheResponse{cacheRepository, nil, c}
}

func runPublisherAsync(resp JetsResponse, checker utils.Checker) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ch := make(chan error, 1)
	go func(c utils.Checker) {
		ch <- ExecutePublisher(resp.Js, resp.Conn, c)
		close(ch)
	}(checker)

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("publisher error: %v", err)
		}
	case <-ctx.Done():
		go func() {
			if err := <-ch; err != nil {
				log.Printf("Delayed publisher error: %v", err)
				return
			}
		}()
	}
	return nil
}

func startServer(cfg *config.Config, controller *chi.Mux) error {
	customServer := server.NewServer(controller,
		server.SetReadTimeout(6*time.Second),
		server.SetWriteTimeout(6*time.Second),
		server.SetAddr(),
		server.SetShutdownTimeout(cfg.Server.ShutdownTimeout),
	)
	if err := customServer.GracefulShutdown(); err != nil {
		return fmt.Errorf("server shutdown error: %v", err)
	}
	return nil
}

func RunApp(cfg *config.Config) error {
	db, err := initPostgres(cfg)
	if err != nil {
		return err
	}
	defer func(db *postgres.DatabaseSource) {
		db.Close()
	}(db)

	response := InitJetStream(cfg)
	if response.Err != nil {
		return response.Err
	}
	natsSrc, err := internal_nats.NewNatsSource(response.Conn, response.Js)
	if err != nil {
		return fmt.Errorf("nats source error: %v", err)
	}

	defer func() {
		if err = natsSrc.Close(); err != nil {
			log.Printf("Nats source close error: %v\n", err)
		}
	}()

	checker := make(map[string]struct{})
	pgRepository := initDBRepository(db)
	cacheResp := initCache(pgRepository, checker)
	if cacheResp.Err != nil {
		return fmt.Errorf("error during creation of repository: %v\n", cacheResp.Err)
	}
	natsRepository := natsstreaming.NewNatsJetStreamRepository(natsSrc)

	natsService := service.CreateNewNatsService(pgRepository, *natsRepository, cacheResp.Repo)
	subscription, err := natsService.StartSubscribing(Channel, QueueGroup)
	defer func() {
		if err = subscription.Unsubscribe(); err != nil {
			log.Printf("Error unsubscribing: %v\n", err)
		}
	}()
	if err != nil {
		return fmt.Errorf("failed to start nats service: %v", err)
	}

	orderService := service.CreateNewOrderService(cacheResp.Repo)
	orderController := controller.CreateNewOrderController(orderService)
	if err = runPublisherAsync(response, checker); err != nil {
		return fmt.Errorf("runPublisherAsync error: %v", err)
	}
	if err = startServer(cfg, orderController); err != nil {
		return fmt.Errorf("start server error: %v", err)
	}
	return nil
}
