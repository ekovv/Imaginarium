package redis

import (
	"Imaginarium/config"
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type DRedisStorage struct {
	conn *redis.Client
}

func NewRedisStorage(config config.Config) *DRedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // адрес сервера Redis
		Password: "",               // пароль (если есть)
		DB:       0,                // используемая база данных
	})
	s := &DRedisStorage{
		conn: client,
	}
	return s
}

func (s *DRedisStorage) Close() error {
	return s.conn.Close()
}

func (s *DRedisStorage) Save(name string, id int) error {
	ctx := context.TODO()
	err := s.conn.Set(ctx, name, id, 0).Err()
	if err != nil {
		return fmt.Errorf("not save in database: %w", err)
	}
	return nil
}

func (s *DRedisStorage) TakeNickName(id int) (string, error) {
	panic("implement me")
}

func (s *DRedisStorage) TakeID(nick string) (int, error) {
	panic("implement me")
}

func (s *DRedisStorage) SavePoints(idOfUser int, nickName string, points int, chatID int) error {
	panic("implement me")
}

func (s *DRedisStorage) TakeAllPoints(chatID int) ([][]string, error) {
	panic("implement me")
}

// doing
func (s *DRedisStorage) SetFall(chatID int, people ...string) error {
	panic("implement me")
}
