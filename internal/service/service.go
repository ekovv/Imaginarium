package service

import (
	"Imaginarium/internal/shema"
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
	AddInMap(chatID int, userID int) (map[int][]shema.Gamers, error)
	Association(association string, userID int, chatID int) (string, int, error)
	MapIsFull(chatID int, userID int) bool
	StartG(chatID int) (string, error)
	TakePhoto(userID int, photoNumber int) (int, []shema.Gamers, error)
	Vote(vote int, userID int, chatID int) ([]shema.Voting, *tele.Photo, error)
	Logic(vote []shema.Voting, chatID int)
}

type Service struct {
	Storage          storage.Storage
	game             map[int][]shema.Gamers
	wantPlay         map[int][]int
	countCards       map[int]int
	countPlayers     map[int]int
	countAssociation map[int]int
	countReady       map[int]int
	flag             bool
	inGame           map[int][]shema.Gamers
	voting           map[int][]shema.Voting
	IdOfAssociated   map[int]int
	resultOfVoting   map[string]int
}

func NewService(storage storage.Storage) *Service {
	return &Service{
		Storage: storage, game: make(map[int][]shema.Gamers),
		wantPlay: make(map[int][]int), inGame: make(map[int][]shema.Gamers),
		voting:           make(map[int][]shema.Voting),
		resultOfVoting:   make(map[string]int),
		countCards:       make(map[int]int),
		countPlayers:     make(map[int]int),
		countAssociation: make(map[int]int),
		countReady:       make(map[int]int),
	}
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
	s.countPlayers[chatID]++
	s.countCards[chatID]++
	s.wantPlay[chatID] = append(s.wantPlay[chatID], userID)
	return nil
}

func (s *Service) AddInMap(chatID int, userID int) (map[int][]shema.Gamers, error) {
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
	g := shema.Gamers{}
continiueLoop:
	for _, file := range files {
		p, _ := s.game[chatID]
		for _, i := range p {
			if i.ID == userID && len(i.Img) >= s.countPlayers[chatID] {
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
		if len(g.Img) != s.countCards[chatID] {
			continue
		}
		s.game[chatID] = append(s.game[chatID], g)

	}
	return s.game, nil
}

func (s *Service) Association(association string, userID int, chatID int) (string, int, error) {
	for key, value := range s.wantPlay {
		for _, user := range value {
			if user == userID {
				result := strings.TrimPrefix(association, "/")
				s.countAssociation[chatID]++
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
					s.countReady[chatID]++
					s.flag = true
				} else {
					s.flag = false
				}
			}
		}
	}
	if s.flag && s.countReady[chatID] == s.countPlayers[chatID] {
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
	for i, d := range array {
		if index == i {
			nickName, err := s.Storage.TakeNickName(array[i])
			if err != nil {
				return "", err
			}
			s.IdOfAssociated[chatID] = d
			return nickName, nil
		}
	}
	return "", nil

}

func (s *Service) TakePhoto(userID int, photoNumber int) (int, []shema.Gamers, error) {
	chatID := 0
	for k, v := range s.game {
		for _, d := range v {
			if d.ID == userID {
				for i, p := range d.Img {
					if i == photoNumber {
						chatID = k
						gamer := shema.Gamers{}
						gamer.ID = userID
						gamer.Img = append(gamer.Img, p)
						s.inGame[chatID] = append(s.inGame[chatID], gamer)
					}
				}
			}
		}
	}
	for _, o := range s.inGame {
		if len(o) == s.countPlayers[chatID] {
			for _, x := range o {
				if len(x.Img) == 1 && chatID != 0 {
					return chatID, o, nil
				}
			}
		}
	}
	return 0, nil, nil
}

func (s *Service) Vote(vote int, userID int, chatID int) ([]shema.Voting, *tele.Photo, error) {
	userWinID := 0
	for k, v := range s.inGame {
		if k == chatID {
			for _, x := range v {
				if x.ID != s.IdOfAssociated[chatID] {
					for i, d := range x.Img {
						for _, j := range s.game {
							for _, q := range j {
								for _, a := range q.Img {
									if a == d {
										userWinID = q.ID
									}
								}
							}
						}
						if x.ID == userWinID && vote == i {
							for g, h := range s.voting {
								if g == chatID {
									for _, y := range h {
										if y.IDWin == userWinID {
											nickNameVote, err := s.Storage.TakeNickName(userID)
											if err != nil {
												return nil, nil, err
											}
											y.NicknameVote = append(y.NicknameVote, nickNameVote)
											y.Count++
										} else {
											vot := shema.Voting{}
											vot.IDWin = userWinID
											nickNameWin, err := s.Storage.TakeNickName(userWinID)
											nickNameVote, err := s.Storage.TakeNickName(userID)
											if err != nil {
												return nil, nil, err
											}
											vot.NicknameWin = "@" + nickNameWin
											vot.NicknameVote = append(vot.NicknameVote, "@"+nickNameVote)
											vot.Count++
											s.voting[chatID] = append(s.voting[chatID], vot)
										}

									}
								}
							}
						}
						if len(s.voting[chatID]) == s.countPlayers[chatID]-1 {
							var photoWin *tele.Photo
							for n, l := range s.inGame {
								if n == chatID {
									for _, m := range l {
										if m.ID == userWinID {
											for _, z := range m.Img {
												photoWin = z
											}
										}
									}
								}
							}
							return s.voting[chatID], photoWin, nil
						}
					}
				}
			}
		}
	}
	return nil, nil, nil
}

func (s *Service) Logic(vote []shema.Voting, chatID int) ([]shema.Points, error) {
	for _, v := range vote {
		if v.IDWin == s.IdOfAssociated[chatID] {
			if v.NicknameVote != nil {
				if len(v.NicknameVote) == s.countPlayers[chatID]-1 {
					s.resultOfVoting[v.NicknameWin] -= 3
					continue
				}
			}
			if v.NicknameVote == nil {
				s.resultOfVoting[v.NicknameWin] -= 2
				continue
			} else {
				for _, n := range v.NicknameVote {
					s.resultOfVoting[n] += 3
				}
			}
			s.resultOfVoting[v.NicknameWin] += 3 + v.Count
		} else {
			s.resultOfVoting[v.NicknameWin] += v.Count
		}

	}
	return nil, nil
}
