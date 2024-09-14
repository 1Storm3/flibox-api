package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(connectionString string) (*Storage, error) {

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL database!")

	return &Storage{db: db}, nil
}
func (s *Storage) Close() error {
	return s.db.Close()
}
