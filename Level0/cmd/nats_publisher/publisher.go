package main

import (
	"Level0/internal/entity"
	"Level0/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/nats-io/stan.go"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	MinTime = 10 * time.Millisecond
	MaxTime = 5000 * time.Millisecond
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// не хардкодить .json файл
func getGivenOrder() entity.Order {
	file, err := os.Open("model.json")
	if err != nil {
		log.Fatal(err)
		return entity.Order{}
	}
	var order entity.Order
	err = json.NewDecoder(file).Decode(&order)
	if err != nil {
		log.Fatal(err)
		return entity.Order{}
	}
	return order
}

// PUBLISHER, основная "труба"
// подумать над генерацией идентификаторов клиентов и имён кластеров в NATS и их хранением
// пока что в .env
func main() {
	sc, err := stan.Connect(os.Getenv("CLUSTER_ID"), os.Getenv("CLIENT_ID"), stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal(fmt.Sprintf("%v", err))
	}
	// получаем ссылку на последний заказ из  model.json
	lastOrder := getGivenOrder()

	// не хардкодить канал
	channel := "subject"
	// подумать над использованием контекста
	for counter := 1; ; counter++ {
		marshalled, err := json.Marshal(lastOrder)
		if err != nil {
			log.Fatal(err)
			return
		}
		err = sc.Publish(channel, marshalled)
		if err != nil {
			log.Fatal(err)
			return
		}
		lastOrder = utils.RandomOrder()
		fmt.Printf("Published Order: %s on channel %s\n", marshalled, channel)
		jitter := func(min, max time.Duration) time.Duration {
			if min >= max {
				return min
			}
			delta := max - min
			return min + time.Duration(rand.Int63n(int64(delta)))
		}(MinTime, MaxTime)
		time.Sleep(jitter)
	}
}
