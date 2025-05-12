package utils

import (
	"Level0/internal/entity"
	"encoding/json"
	"log"
	"os"
)

// не хардкодить .json файл
func GetGivenOrder() entity.Order {
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
