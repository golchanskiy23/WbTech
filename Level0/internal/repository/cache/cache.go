package cache

import (
	"Level0/internal/entity"
	DBRepo "Level0/internal/repository/database"
	"context"
	"fmt"
	"sync"
)

type CacheRepository struct {
	Mtx   *sync.RWMutex
	Cache map[string]*entity.Order
}

func CreateNewCacheRepository(pg *DBRepo.DatabaseRepository) (*CacheRepository, error) {
	cacheRepository := &CacheRepository{
		Mtx:   &sync.RWMutex{},
		Cache: make(map[string]*entity.Order),
	}
	orders, err := pg.GetAllOrders(context.Background())
	if err != nil {
		fmt.Println()
		return nil, err
	}
	for _, order := range orders {
		cacheRepository.Mtx.Lock()
		cacheRepository.Cache[order.OrderUID] = &order
		cacheRepository.Mtx.Unlock()
	}
	return cacheRepository, nil
}
