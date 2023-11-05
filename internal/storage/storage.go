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

func (s *Storage) Save(name string, id int) error {
	insertQuery := "INSERT INTO users(name, idoftelegram) VALUES ($1, $2)"
	_, err := s.conn.Exec(insertQuery, name, id)
	if err != nil {
		return fmt.Errorf("not save in database: %w", err)
	}
	return nil
}

func (s *Storage) TakeNickName(id int) (string, error) {
	query := "SELECT name FROM users WHERE idoftelegram = $1"
	var nickName string
	err := s.conn.QueryRow(query, id).Scan(&nickName)
	if err != nil {
		return "", fmt.Errorf("not take name from database: %w", err)
	}
	return nickName, nil
}
