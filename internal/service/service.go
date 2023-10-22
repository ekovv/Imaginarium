package service

import (
	"Imaginarium/internal/storage"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"os"
)

type Inter interface {
	SaveInDB(id int) error
	Inc(chatID int, userID int) error
	AddInMap(chatID int, userID int) (map[int]Gamers, error)
}

type Service struct {
	Storage      storage.Storage
	game         map[int]Gamers
	wantPlay     map[int][]int
	countCards   int
	countPlayers int
}

func NewService(storage storage.Storage) *Service {
	return &Service{Storage: storage, game: make(map[int]Gamers), wantPlay: make(map[int][]int)}
}

func (s *Service) SaveInDB(id int) error {
	err := s.Storage.Save(id)
	if err != nil {
		return fmt.Errorf("not save in database: %w", err)
	}
	return nil
}

func (s *Service) Inc(chatID int, userID int) error {
	_, ok := s.game[chatID]
	if ok {
		return fmt.Errorf("Game in process")
	}
	s.countPlayers++
	s.countCards++
	s.wantPlay[chatID] = append(s.wantPlay[chatID], userID)

	return nil
}

func (s *Service) AddInMap(chatID int, userID int) (map[int]Gamers, error) {
	_, ok := s.game[chatID]
	if ok {
		return nil, fmt.Errorf("Game in process")
	}
	files, err := os.ReadDir("./src")
	if err != nil {
		fmt.Println("Ошибка чтения папки:", err)
		return nil, err
	}
	for _, e := range s.wantPlay[chatID] {
		for _, file := range files {
			if len(s.game[e].Img) >= s.countCards {
				break
			}
			if file.Name() == ".DS_Store" {
				continue
			}
			photo := &tele.Photo{File: tele.FromDisk(file.Name())}
			g := Gamers{}
			g.ID = userID
			g.Img = append(g.Img, photo)
			s.game[chatID] = g

		}
	}
	return s.game, nil
}
