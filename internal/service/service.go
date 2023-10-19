package service

import (
	"Imaginarium/internal/storage"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"os"
	"path/filepath"
	"strings"
)

type Inter interface {
	SaveInDB(id int) error
	Inc(id int) error
	AddInMap() map[int][]*tele.Photo
}

type Service struct {
	Storage      storage.Storage
	game         map[int][]*tele.Photo
	wantPlay     []int
	countCards   int
	countPlayers int
}

func NewService(storage storage.Storage) *Service {
	return &Service{Storage: storage, game: make(map[int][]*tele.Photo)}
}

func (s *Service) SaveInDB(id int) error {
	err := s.Storage.Save(id)
	if err != nil {
		return fmt.Errorf("not save in database: %w", err)
	}
	return nil
}

func (s *Service) Inc(id int) error {
	s.countPlayers++
	s.countCards++
	s.wantPlay = append(s.wantPlay, id)
	return nil
}

func (s *Service) AddInMap() map[int][]*tele.Photo {
	files, err := os.ReadDir("./src")
	if err != nil {
		fmt.Println("Ошибка чтения папки:", err)
		return nil
	}
	for i, _ := range s.wantPlay {
		for _, file := range files {
			if isImage(file.Name()) {
				photo := &tele.Photo{File: tele.FromDisk(file.Name())}
				s.game[i] = append(s.game[i], photo)
				if len(s.game[i]) >= s.countCards {
					break
				}
			}
		}
	}
	return s.game
}

func isImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png"
}
