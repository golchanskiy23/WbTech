package service

import (
	"Level0/internal/entity"
	"Level0/internal/repository/cache"
	"Level0/internal/repository/database"
	"Level0/internal/repository/natsstreaming"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"log"
)

type NatsStreamingService interface {
	StartSubscribing(channel, queue string) (stan.Subscription, error)
	handleMessage(msg *stan.Msg) error
}

type NatsService struct {
	PgRepository    database.CRUDRepository
	NatsRepository  natsstreaming.NatsJetStreamRepository
	CacheRepository *cache.CacheRepository
}

func CreateNewNatsService(pg database.CRUDRepository, nats natsstreaming.NatsJetStreamRepository, cache *cache.CacheRepository) NatsService {
	return NatsService{PgRepository: pg, NatsRepository: nats, CacheRepository: cache}
}

func (s NatsService) AddOrderToDB(order entity.Order) {
	err := s.PgRepository.AddOrder(context.Background(), order)
	if err != nil {
		log.Fatalf("Impossible to add order to DB #665: %v", err)
		return
	}
}

func (s NatsService) AddOrderToCache(order *entity.Order) {
	s.CacheRepository.Set(order)
	for k, _ := range s.CacheRepository.Cache {
		fmt.Printf("%d ", k)
	}
	fmt.Println()
}

func (service NatsService) StartSubscribing(subject, durable string) (*nats.Subscription, error) {
	return service.NatsRepository.NatsSrc.JetStream.QueueSubscribe(subject, durable, func(msg *nats.Msg) {
		if err := service.handleMessage(msg); err != nil {
			log.Fatal(err)
			return
		}
		var order entity.Order
		json.Unmarshal(msg.Data, &order)
		service.AddOrderToDB(order)
		service.AddOrderToCache(&order)
		if err := msg.Ack(); err != nil {
			log.Printf("ack error: %v", err)
		}
	},
		nats.Durable(durable),
	)
}

func (service NatsService) handleMessage(msg *nats.Msg) error {
	var order entity.Order
	err := json.Unmarshal(msg.Data, &order)
	if err != nil {
		return err
	}
	service.CacheRepository.Set(&order)
	return nil
}
