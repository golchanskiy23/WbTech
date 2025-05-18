package service

import (
	"Level0/internal/entity"
	"Level0/internal/repository/cache"
	"Level0/internal/repository/database"
	"Level0/internal/repository/natsstreaming"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
)

type NatsService struct {
	PgRepository    *database.DatabaseRepository
	NatsRepository  *natsstreaming.NatsStreamingRepository
	CacheRepository *cache.CacheRepository
}

func CreateNewNatsService(pg *database.DatabaseRepository, nats *natsstreaming.NatsStreamingRepository, cache *cache.CacheRepository) *NatsService {
	return &NatsService{PgRepository: pg, NatsRepository: nats, CacheRepository: cache}
}

func (s *NatsService) AddOrderToDB(order entity.Order) {
	err := s.PgRepository.AddOrder(context.Background(), order)
	if err != nil {
		log.Fatalf("Impossible to add order to DB #665: %v", err)
		return
	}
}

func (s *NatsService) AddOrderToCache(order *entity.Order) {
	s.CacheRepository.Set(order)
	for k, _ := range s.CacheRepository.Cache {
		fmt.Printf("%d ", k)
	}
	fmt.Println()
}

func (service *NatsService) StartSubscribing(channel, queue_group string) (stan.Subscription, error) {
	return service.NatsRepository.NatsSrc.Conn.QueueSubscribe(channel, queue_group, func(msg *stan.Msg) {
		if err := service.handleMessage(msg); err != nil {
			log.Fatal(err)
			return
		}
		var order entity.Order
		json.Unmarshal(msg.Data, &order)
		service.AddOrderToDB(order)
		service.AddOrderToCache(&order)
		//fmt.Println("Received a message:", string(msg.Data))
	})
}

func (service *NatsService) handleMessage(msg *stan.Msg) error {
	var order entity.Order
	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		return err
	}
	service.CacheRepository.Set(&order)
	return nil
}
