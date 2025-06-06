package main

import "errors"

var (
	OK     = "TRANSACTION_COMPLETED"
	NOT_OK = "TRANSACTION_UNCOMPLETED"
)

type NotificationService struct {
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (service *NotificationService) Notify(before, after int) (string, error) {
	if before > after {
		return OK, nil
	}
	return NOT_OK, errors.New("Notification service not available")
}
