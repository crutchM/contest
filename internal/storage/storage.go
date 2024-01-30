package storage

import (
	"Contest/internal/domain"
	"Contest/internal/storage/postgres"
	"database/sql"
	_ "github.com/lib/pq"
)

type Repository[T any] interface {
	AddItem(item T) error
	DeleteItem(id int) error
	UpdateItem(id int, newItem T) error
	GetTable() ([]T, error)
	FindItemByID(id int) (T, error)
	FindItemByCondition(condition func(item T) bool) (T, error)
	FindItemsByCondition(condition func(item T) bool) ([]T, error)
}

type Storage struct {
	db             *sql.DB
	TestRepository Repository[domain.Test]
}

func NewStorage(connStr string) (*Storage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	testRepository := postgres.NewTestRepository(db)

	return &Storage{
		db:             db,
		TestRepository: testRepository,
	}, nil
}
