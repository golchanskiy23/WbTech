package database

import (
	"Level0/internal/entity"
	"Level0/pkg/postgres"
	"context"
)

type DatabaseRepository struct {
	DB *postgres.DatabaseSource
}

// скрывает базу данных
func CreateNewDBRepository(db *postgres.DatabaseSource) *DatabaseRepository {
	return &DatabaseRepository{DB: db}
}

// набор методов, реализующих работу с базами данных
func (r *DatabaseRepository) GetAllOrders(ctx context.Context) ([]entity.Order, error) {
	return nil, nil
}
