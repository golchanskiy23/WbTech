package main

import (
	"errors"
	"fmt"
)

type FacadeWallet struct {
	Wallet              *Wallet
	Account             *Account
	SecurityCodeService *SecurityCodeService
	NotificationService *NotificationService
}

func NewFacadeWallet(stringID, code string) *FacadeWallet {
	return &FacadeWallet{
		Wallet:              NewWallet(),
		Account:             NewAccount(stringID),
		SecurityCodeService: NewSecurityCodeService(code),
		NotificationService: NewNotificationService(),
	}
}

func CheckAuthentication(f *FacadeWallet, stringID, code string) error {
	if err := f.Account.CheckAccount(stringID); err != nil {
		return errors.New("check account failed")
	}
	if err := f.SecurityCodeService.CheckCode(code); err != nil {
		return errors.New("check code failed")
	}
	return nil
}

func (f *FacadeWallet) AddMoney(stringID, code string, money int) error {
	if err := CheckAuthentication(f, stringID, code); err != nil {
		return err
	}
	before := f.Wallet.GetBalance()
	if err := f.Wallet.Add(money); err != nil {
		return fmt.Errorf("addition money error: %v", err)
	}
	notification, err := f.NotificationService.Notify(before, f.Wallet.GetBalance())
	if err != nil {
		return fmt.Errorf("addition money error in notification: %v", err)
	}
	fmt.Println(notification)
	return nil
}

func (f *FacadeWallet) PayOrder(stringID, code string, order *Order) error {
	if err := CheckAuthentication(f, stringID, code); err != nil {
		return fmt.Errorf("authentication error : %v", err)
	}
	sum, before := 0, f.Wallet.GetBalance()
	for _, o := range order.Items {
		sum += o.Value
	}
	if f.Wallet.GetBalance() < sum {
		return errors.New("not enough balance")
	}
	if err := f.Wallet.Subtract(sum); err != nil {
		return errors.New("not enough balance during subtraction")
	}
	notification, err := f.NotificationService.Notify(before, f.Wallet.GetBalance())
	if err != nil {
		return errors.New("balance isn't decreasing")
	}
	fmt.Println(notification)
	return nil
}
