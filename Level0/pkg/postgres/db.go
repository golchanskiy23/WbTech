package postgres

import (
	"Level0/config"
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
		os.Getenv("POSTRGES_UNSAFE_USERNAME"),
		os.Getenv("POSTRGES_UNSAFE_PASSWORD"),
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)
	return dsn
}

func (s *DatabaseSource) Close() {
	s.Pool.Close()
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
		return nil, err
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
			return nil, errors.New("max connection attempts exceeded")
		}

		jitter := utils.CreateNewDelay(attempt, src.MaxConnectTimeout)
		fmt.Printf("Attempt number %d failed, waiting for %v\n", attempt+1, jitter)
		time.Sleep(jitter)
	}
	return nil, ctx.Err()
}
