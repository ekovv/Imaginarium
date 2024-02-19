package storage

import (
	"Imaginarium/config"
	"Imaginarium/internal/storage/postgresql"
	"Imaginarium/internal/storage/redis"
	"fmt"
)

type Storage interface {
	SetFall(chatID int, people ...string) error
	TakeAllPoints(chatID int) ([][]string, error)
	SavePoints(idOfUser int, nickName string, points int, chatID int) error
	TakeID(nick string) (int, error)
	TakeNickName(id int) (string, error)
	Save(name string, id int) error
}

// New sa
func New(cfg config.Config) (Storage, error) {
	switch cfg.TypeDB {
	case "postgresql":
		d := postgresql.NewPostgresqlStorage(cfg)
		return d, nil
	case "redis":
		r := redis.NewRedisStorage(cfg)
		return r, nil
	}
	return nil, fmt.Errorf("unsupported storage type")
}
