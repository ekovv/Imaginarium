package storage

import (
	"Imaginarium/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	conn *sql.DB
}

func NewStorage(config config.Config) *Storage {
	db, err := sql.Open("postgres", config.DB)
	if err != nil {
		fmt.Println("Error connect to db")
	}
	s := &Storage{
		conn: db,
	}
	return s
}

func (s *Storage) Close() error {
	return s.conn.Close()
}

func (s *Storage) Save(name string) error {
	insertQuery := "INSERT INTO users(name) VALUES ($1)"
	_, err := s.conn.Exec(insertQuery, name)
	if err != nil {
		return fmt.Errorf("not save in database: %w", err)
	}
	return nil
}
