package service

import (
	"Imaginarium/internal/storage"
	"fmt"
)

type Inter interface {
	SaveInDB(id int) error
	AddInMap(id int) error
}

type Service struct {
	Storage storage.Storage
	game    map[int][]string
}

func NewService(storage storage.Storage) *Service {
	return &Service{Storage: storage, game: make(map[int][]string)}
}

func (s *Service) SaveInDB(id int) error {
	err := s.Storage.Save(id)
	if err != nil {
		return fmt.Errorf("not save in database: %w", err)
	}
	return nil
}

func (s *Service) AddInMap(id int) error {
	_, ok := s.game[id]
	if !ok {
		s.game[id] = []string{"123"}
	}
	return nil
}
