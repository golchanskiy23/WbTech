package main

import (
	"fmt"
	"log"
)

func main() {
	facade := NewFacadeWallet("345ABC", "123")

	if err := facade.AddMoney("345ABC", "123", 1800); err != nil {
		fmt.Println("Try again, error during replenishment of the wallet")
	}
	fmt.Println("The replenishment of the wallet was successful")
	order := &Order{
		Items: []Item{{"first", 546}, {"second", 547}, {"third", 500}},
	}
	if err := facade.PayOrder("345ABC", "123", order); err != nil {
		log.Fatalf("Fatal error during payment: %v", err)
	}
	fmt.Println("Successful payment!!!")
}
