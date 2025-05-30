package cache

import (
	"Level0/internal/entity"
	"errors"
	"log"
	"sync"
)

type CacheRepository struct {
	Mtx   *sync.RWMutex
	Cache map[string]*entity.Order
}

func (repo CacheRepository) IsEmpty() bool {
	return repo.Cache == nil
}

func CreateNewCacheRepository(orders []entity.Order) (*CacheRepository, error) {
	cacheRepository := &CacheRepository{
		Mtx:   &sync.RWMutex{},
		Cache: make(map[string]*entity.Order),
	}
	log.Printf("Got %d orders from DB", len(orders))
	for _, order := range orders {
		// fmt.Print(order, " ")
		// fmt.Println()
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
