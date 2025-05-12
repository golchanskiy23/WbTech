package database

import "Level0/pkg/postgres"

type DatabaseRepository struct {
	DB *postgres.DatabaseSource
}

// скрывает базу данных
func CreateNewDBRepository(db *postgres.DatabaseSource) *DatabaseRepository {
	return &DatabaseRepository{DB: db}
}

// набор методов, реализующих работу с базами данных
