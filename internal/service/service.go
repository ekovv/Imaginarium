package service

import (
	"Imaginarium/internal/storage"
	"fmt"
)

type Inter interface {
	SaveInDB(name string) error
}

type Service struct {
	Storage storage.Storage
	m       map[string]interface{}
}

func NewService(storage storage.Storage) *Service {
	return &Service{Storage: storage, m: make(map[string]interface{})}
}

func (s *Service) SaveInDB(name string) error {
	err := s.Storage.Save(name)
	if err != nil {
		return fmt.Errorf("not save in database: %w", err)
	}
	return nil
}
