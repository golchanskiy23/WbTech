package main

import (
	"errors"
	"fmt"
)

type Wallet struct {
	balance int
}

func NewWallet() *Wallet {
	return &Wallet{}
}

func (w *Wallet) Add(amount int) error {
	w.balance += amount
	fmt.Printf("Added %d to the wallet\n", w.balance)
	return nil
}

func (w Wallet) GetBalance() int {
	return w.balance
}

func (w *Wallet) Subtract(money int) error {
	if w.balance < money {
		return errors.New("Not enough money")
	}
	w.balance -= money
	return nil
}
