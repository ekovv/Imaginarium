package postgresql

import (
	"Imaginarium/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type PostgreStorage struct {
	conn *sql.DB
}

func NewPostgresqlStorage(config config.Config) *PostgreStorage {
	db, err := sql.Open("postgres", config.DB)
	if err != nil {
		fmt.Println("Error connect to db")
	}
	s := &PostgreStorage{
		conn: db,
	}
	return s
}

func (s *PostgreStorage) Close() error {
	return s.conn.Close()
}

func (s *PostgreStorage) Save(name string, id int) error {
	insertQuery := "INSERT INTO users(name, idoftelegram) VALUES ($1, $2)"
	_, err := s.conn.Exec(insertQuery, name, id)

	if err != nil {
		return fmt.Errorf("not save in database: %w", err)
	}
	return nil
}

func (s *PostgreStorage) TakeNickName(id int) (string, error) {
	query := "SELECT name FROM users WHERE idoftelegram = $1"
	var nickName string
	err := s.conn.QueryRow(query, id).Scan(&nickName)
	if err != nil {
		return "", fmt.Errorf("not take name from database: %w", err)
	}
	return nickName, nil
}

func (s *PostgreStorage) TakeID(nick string) (int, error) {
	query := "SELECT idoftelegram FROM users WHERE name = $1"
	var idOfTelegram int
	err := s.conn.QueryRow(query, nick).Scan(&idOfTelegram)
	if err != nil {
		return 0, fmt.Errorf("not take name from database: %w", err)
	}
	return idOfTelegram, nil
}

func (s *PostgreStorage) SavePoints(idOfUser int, nickName string, points int, chatID int) error {
	insertQuery := "INSERT INTO points(user_id, nickname, score, chat, flag) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (user_id) DO UPDATE set score = $6"
	_, err := s.conn.Exec(insertQuery, idOfUser, nickName, points, chatID, true, points)
	if err != nil {
		return fmt.Errorf("not save in database: %w", err)
	}
	return nil
}

func (s *PostgreStorage) TakeAllPoints(chatID int) ([][]string, error) {
	rows, err := s.conn.Query("SELECT chat,nickname,score FROM points WHERE chat = $1 AND flag = true", chatID)
	if err != nil {
		log.Fatal(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var result [][]string
	cols, _ := rows.Columns()

	for rows.Next() {
		columns := make([]string, len(cols))
		colPtrs := make([]interface{}, len(cols))
		for i := range columns {
			colPtrs[i] = &columns[i]
		}

		if err := rows.Scan(colPtrs...); err != nil {
			log.Fatal(err)
		}

		result = append(result, columns)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return result, nil
}

// doing
func (s *PostgreStorage) SetFall(chatID int, people ...string) error {
	return nil
}
