package main

import (
	"errors"
	"fmt"
)

type SecurityCodeService struct {
	Code string
}

func NewSecurityCodeService(code string) *SecurityCodeService {
	return &SecurityCodeService{Code: code}
}

func (service *SecurityCodeService) CheckCode(s string) error {
	if service.Code != s {
		return errors.New("Code does not match")
	}
	fmt.Println("Security Code Found")
	return nil
}
