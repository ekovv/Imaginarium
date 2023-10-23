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
	AddInMap(chatID int, userID int) (map[int][]Gamers, error)
}

type Service struct {
	Storage      storage.Storage
	game         map[int][]Gamers
	wantPlay     map[int][]int
	countCards   int
	countPlayers int
}

func NewService(storage storage.Storage) *Service {
	return &Service{Storage: storage, game: make(map[int][]Gamers), wantPlay: make(map[int][]int)}
}

func (s *Service) SaveInDB(id int) error {
	err := s.Storage.Save(id)
	if err != nil {
		return fmt.Errorf("not save in database: %w", err)
	}
	return nil
}

func (s *Service) Inc(chatID int, userID int) error {
	ph, _ := s.game[chatID]
	for _, i := range ph {
		if i.Img != nil && i.ID == userID {
			return fmt.Errorf("Game in process")
		}
	}
	s.countPlayers++
	s.countCards++
	s.wantPlay[chatID] = append(s.wantPlay[chatID], userID)
	return nil
}

func (s *Service) AddInMap(chatID int, userID int) (map[int][]Gamers, error) {
	ph, _ := s.game[chatID]
	for _, i := range ph {
		if i.Img != nil && i.ID == userID {
			return nil, fmt.Errorf("Game in process")
		}
	}
	files, err := os.ReadDir("./src")
	if err != nil {
		fmt.Println("Ошибка чтения папки:", err)
		return nil, err
	}
	g := Gamers{}
	for _, file := range files {
		p, _ := s.game[chatID]
		for _, i := range p {
			if i.ID == userID && len(i.Img) >= s.countPlayers {
				return s.game, nil
			}
		}
		if file.Name() == ".DS_Store" {
			continue
		}
		photo := &tele.Photo{File: tele.FromDisk(file.Name())}
		g.ID = userID
		g.Img = append(g.Img, photo)
		if len(g.Img) != s.countCards {
			continue
		}
		s.game[chatID] = append(s.game[chatID], g)

	}
	return s.game, nil
}
