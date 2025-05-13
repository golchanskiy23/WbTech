package app

import (
	"Level0/config"
	"Level0/internal/repository/cache"
	"Level0/internal/repository/database"
	"Level0/internal/repository/natsstreaming"
	"Level0/internal/service"
	"Level0/internal/utils"
	"Level0/pkg/nats"
	"Level0/pkg/postgres"
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func InitInsertingDB(db *postgres.DatabaseSource) error {
	order := utils.GetGivenOrder()
	_, err := db.Pool.Exec(context.Background(), `
	INSERT INTO orders (
		order_uid, track_number, entry, locale, internal_signature,
		customer_id, delivery_service, shard_key, sm_id, data_created, oof_shard
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DataCreated, order.OofShard)
	if err != nil {
		log.Fatal("insert orders:", err)
		return err
	}

	_, err = db.Pool.Exec(context.Background(), `
	INSERT INTO deliveries (
		order_uid, name, phone, zip, city, address, region, email
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
`, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		log.Fatal("insert deliveries:", err)
		return err
	}

	_, err = db.Pool.Exec(context.Background(), `
	INSERT INTO payments (
		order_uid, transaction, request_id, currency, provider,
		amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
`, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT,
		order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		log.Fatal("insert payments:", err)
		return err
	}

	for _, item := range order.Items {
		_, err := db.Pool.Exec(context.Background(), `
		INSERT INTO items (
			order_uid, chrt_id, track_number, price, rid,
			name, sale, size, total_price, nm_id, brand, status
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid,
			item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			log.Println("insert item:", err)
			return err
		}
	}
	return nil
}

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

	err = InitInsertingDB(db)
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

	natsService := service.CreateNewNatsService(pgRepository, natsRepository, cacheRepository)
	subscription, err := natsService.StartSubscribing("subject", "queue_group")
	defer subscription.Close()
	if err != nil {
		fmt.Println("4")
		fmt.Printf("Failed to start nats service: %v", err)
	}

	/*
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
