package cache

import (
	"Level0/internal/entity"
	DBRepo "Level0/internal/repository/database"
	"context"
	"errors"
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
		cacheRepository.Set(&order)
	}
	return cacheRepository, nil
}

func (cache *CacheRepository) Set(order *entity.Order) {
	cache.Mtx.Lock()
	defer cache.Mtx.Unlock()
	if _, b := cache.Cache[order.OrderUID]; b {
		return
	}
	cache.Cache[order.OrderUID] = order
}

func (cache *CacheRepository) GetById(id string) (entity.Order, error) {
	if val, isOk := cache.Cache[id]; !isOk {
		return entity.Order{}, errors.New("Order Not Found")
	} else {
		return *val, nil
	}
}
