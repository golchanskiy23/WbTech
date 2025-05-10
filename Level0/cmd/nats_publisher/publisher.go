package main

import (
	"Level0/internal/entity"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	stan "github.com/nats-io/stan.go"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getGivenOrder() *entity.Order {
	return &entity.Order{}
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
	channel := "subject"
	id := lastOrder.UUID
	// публикуем заказ в канал спустя промежутки времени
	for counter := 1; ; counter++ {
		// подумать над случайным заполнением полей заказа
		// currentOrder := CreateNewOrder()
		lastOrder.UUID = id + strconv.Itoa(counter)
		lastOrder.Transaction = lastOrder.UUID
		lastOrder.Items[0] = lastOrder.Items[0] + rand.Int()
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
		fmt.Printf("Published Order: %s on channel %s\n", marshalled, channel)
		time.Sleep(5 * time.Second)
	}
}
