package utils

import (
	"Level0/internal/entity"
	"encoding/json"
	"fmt"
	"os"
)

const (
	ModelFile = "model.json"
)

func GetGivenOrder() (entity.Order, error) {
	file, err := os.Open(ModelFile)
	if err != nil {
		return entity.Order{}, fmt.Errorf("error opening file: %v", err)
	}
	var order entity.Order
	err = json.NewDecoder(file).Decode(&order)
	if err != nil {
		return entity.Order{}, fmt.Errorf("error decoding file: %v", err)
	}
	return order, nil
}
