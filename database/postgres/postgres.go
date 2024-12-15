package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db *gorm.DB
}

func (s *Storage) DB() *gorm.DB {
	return s.db
}

func NewStorage(connectionString string) (*Storage, error) {

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Ошибка при подключении к базе: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("Ошибка при подключении к базе: %v", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("Ошибка при подключении к базе: %v", err)
	}

	fmt.Println("Успешное подключение к базе!")

	return &Storage{db: db}, nil
}
func (s *Storage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
