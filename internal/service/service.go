package service

import (
	"Imaginarium/internal/storage"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Inter interface {
	SaveInDB(name string, id int) error
	Inc(chatID int, userID int) error
	AddInMap(chatID int, userID int) (map[int][]Gamers, error)
	Association(association string, userID int) (string, int, error)
	MapIsFull(chatID int, userID int) bool
	StartG(chatID int) (string, error)
}

type Service struct {
	Storage          storage.Storage
	game             map[int][]Gamers
	wantPlay         map[int][]int
	countCards       int
	countPlayers     int
	countAssociation int
	countReady       int
	flag             bool
}

func NewService(storage storage.Storage) *Service {
	return &Service{Storage: storage, game: make(map[int][]Gamers), wantPlay: make(map[int][]int)}
}

func (s *Service) SaveInDB(name string, id int) error {
	err := s.Storage.Save(name, id)
	if err != nil {
		return fmt.Errorf("not save in database: %w", err)
	}
	return nil
}

func (s *Service) Inc(chatID int, userID int) error {
	ph, _ := s.game[chatID]
	for _, i := range ph {
		if i.Img != nil {
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
continiueLoop:
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
		for _, a := range p {
			for _, o := range a.Img {
				if o.FileLocal == file.Name() {
					continue continiueLoop
				}
			}
		}
		for _, q := range g.Img {
			if q == photo {
				continue
			}
		}
		g.Img = append(g.Img, photo)
		if len(g.Img) != s.countCards {
			continue
		}
		s.game[chatID] = append(s.game[chatID], g)

	}
	return s.game, nil
}

func (s *Service) Association(association string, userID int) (string, int, error) {
	for key, value := range s.wantPlay {
		for _, user := range value {
			if user == userID {
				result := strings.TrimPrefix(association, "/")
				s.countAssociation++
				return result, key, nil
			}
		}
	}
	return "", 0, fmt.Errorf("Нету тебя в беседе")
}

func (s *Service) MapIsFull(chatID int, userID int) bool {
	for key, value := range s.game {
		if value != nil && key == chatID {
			for _, i := range value {
				if i.ID == userID && i.Img != nil {
					s.countReady++
					s.flag = true
				} else {
					s.flag = false
				}
			}
		}
	}
	if s.flag && s.countReady == s.countPlayers {
		return true
	} else {
		return false
	}
}

func (s *Service) StartG(chatID int) (string, error) {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var array []int
	for key, value := range s.game {
		if key == chatID {
			for _, v := range value {
				array = append(array, v.ID)

			}
		}
	}
	index := rand.Intn(len(array))
	for i, _ := range array {
		if index == i {
			nickName, err := s.Storage.TakeNickName(array[i])
			if err != nil {
				return "", err
			}
			return nickName, nil
		}
	}
	return "", nil

}
