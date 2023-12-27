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

func (s *Storage) TakeID(nick string) (int, error) {
	query := "SELECT idoftelegram FROM users WHERE name = $1"
	var idOfTelegram int
	err := s.conn.QueryRow(query, nick).Scan(&idOfTelegram)
	if err != nil {
		return 0, fmt.Errorf("not take name from database: %w", err)
	}
	return idOfTelegram, nil
}

func (s *Storage) SavePoints(idOfUser int, nickName string, points int) error {
	insertQuery := "INSERT INTO points(user_id, nickname, score) VALUES ($1, $2, $3)"
	_, err := s.conn.Exec(insertQuery, idOfUser, nickName, points)
	if err != nil {
		return fmt.Errorf("not save in database: %w", err)
	}
	return nil
}
