package postgres

import (
	"Level0/config"
	"Level0/internal/entity"
	"Level0/internal/utils"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"time"
)

type DatabaseSource struct {
	Pool               *pgxpool.Pool
	MaxPoolSize        int
	MaxConnectTimeout  time.Duration
	MaxConnLifetime    time.Duration
	MaxConnectAttempts int
}

const (
	defaultMaxPoolSize       = 5
	defaultMaxConnLifetime   = 600 * time.Second
	defaultMaxConnectTimeout = 1 * time.Second
	defaultMaxConnAttempts   = 5
)

func GetConnection(cfg *config.DB) string {
	dsn := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		"postgres",
		os.Getenv("POSTGRES_UNSAFE_USERNAME"),
		os.Getenv("POSTGRES_UNSAFE_PASSWORD"),
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)
	// для докера
	/*dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		os.Getenv("POSTGRES_UNSAFE_USERNAME"),
		os.Getenv("POSTGRES_UNSAFE_PASSWORD"),
		cfg.Name,
		cfg.SSLMode,
	)*/
	return dsn
}

func (s *DatabaseSource) Close() {
	s.Pool.Close()
}

func AddOrdersToDB(db *DatabaseSource, order entity.Order) error {
	var exists bool
	sqlString := "SELECT EXISTS(SELECT 1 FROM orders WHERE order_uid = $1)"
	if err := db.Pool.QueryRow(context.Background(), sqlString, order.OrderUID).Scan(&exists); err != nil {
		return fmt.Errorf("error executing sql script for add orders in DB: %w", err)
	}
	if exists {
		return nil
	}
	_, err := db.Pool.Exec(context.Background(), `
	INSERT INTO orders (
		order_uid, track_number, entry, locale, internal_signature,
		customer_id, delivery_service, shard_key, sm_id, data_created, oof_shard
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.ShardKey, order.SmID, order.DataCreated, order.OofShard)
	if err != nil {
		return fmt.Errorf("error inserting sql script for add orders in DB: %w", err)
	}

	_, err = db.Pool.Exec(context.Background(), `
	INSERT INTO deliveries (
		order_uid, name, phone, zip, city, address, region, email
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
`, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return fmt.Errorf("error inserting sql script for add deliveries in DB: %w", err)
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
		return fmt.Errorf("error inserting sql script for add payments in DB: %w", err)
	}

	for _, item := range order.Items {
		_, err = db.Pool.Exec(context.Background(), `
		INSERT INTO items (
			order_uid, chrt_id, track_number, price, rid,
			name, sale, size, total_price, nm_id, brand, status
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
	`, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid,
			item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("error inserting sql script for add items in DB: %w", err)
		}
	}
	return nil
}

func NewStorage(url string, options ...Option) (*DatabaseSource, error) {
	src := &DatabaseSource{
		MaxPoolSize:        defaultMaxPoolSize,
		MaxConnLifetime:    defaultMaxConnLifetime,
		MaxConnectTimeout:  defaultMaxConnectTimeout,
		MaxConnectAttempts: defaultMaxConnAttempts,
	}
	for _, option := range options {
		option(src)
	}
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("creation db storage error: %w", err)
	}
	cfg.MaxConns = int32(src.MaxPoolSize)
	ctx := context.Background()
	for attempt := 0; attempt < src.MaxConnectAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			src.Pool, err = pgxpool.NewWithConfig(context.Background(), cfg)
			if err == nil {
				fmt.Println(src)
				return src, nil
			}
		}
		if attempt == src.MaxConnectAttempts {
			return nil, errors.New("max connection attempts exceeded; connection is failed!")
		}

		jitter := utils.CreateNewDelay(attempt, src.MaxConnectTimeout)
		time.Sleep(jitter)
	}
	return nil, ctx.Err()
}
