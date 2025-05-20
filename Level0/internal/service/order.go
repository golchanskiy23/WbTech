package service

import (
	"Level0/internal/entity"
	"Level0/internal/repository/cache"
	"fmt"
)

type CRUDService interface {
	GetOrderById(id string) (entity entity.Order, err error)
}

type OrderService struct {
	CacheRepo *cache.CacheRepository
}

func (s OrderService) GetOrderById(id string) (entity.Order, error) {
	order, err := s.CacheRepo.GetById(id)
	if err != nil {
		fmt.Errorf("Error getting order by id %s", id)
		return entity.Order{}, err
	}
	return order, nil
}

func CreateNewOrderService(cache *cache.CacheRepository) OrderService {
	return OrderService{
		CacheRepo: cache,
	}
}
