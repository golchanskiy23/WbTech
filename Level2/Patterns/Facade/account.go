package main

import (
	"errors"
	"fmt"
)

type Account struct {
	AccID string
}

func NewAccount(accID string) *Account {
	return &Account{accID}
}

func (account *Account) CheckAccount(id string) error {
	if account.AccID != id {
		return errors.New("Account ID not match")
	}
	fmt.Println("Account check is OK")
	return nil
}
