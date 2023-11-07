package handler

import (
	"Imaginarium/config"
	"Imaginarium/internal/service"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"os"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	Service service.Inter
	Bot     *tele.Bot
}

func NewHandler(service service.Inter, config config.Config) (Handler, error) {
	pref := tele.Settings{
		Token:  config.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return Handler{}, nil
	}
	h := Handler{
		Bot:     b,
		Service: service,
	}
	return h, nil
}

func (s *Handler) Start(c tele.Context) error {
	s.Bot.Send(c.Chat(), "Привет! Я бот, который предлагает вам поиграть в интересную игру.")
	startGameKeyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{
				tele.InlineButton{
					Text: "Начать игру",
					Data: "start_game",
				},
			},
		},
	}
	s.Bot.Reply(c.Message(), "Чтобы начать игру, нажмите на кнопку ниже.", startGameKeyboard)
	return nil
}

func (s *Handler) HandleButton(c tele.Context) error {
	data := c.Data()
	switch data {
	case "start_game":
		s.AddPlayer(c)
	case "\fready":
		s.GiveCards(c)
	case "\f0":
		s.PhotoTake(c)
	case "\f1":
		s.PhotoTake(c)
	case "\f2":
		s.PhotoTake(c)
	case "\f3":
		s.PhotoTake(c)
	case "\f4":
		s.PhotoTake(c)
	}
	return nil
}

func (s *Handler) AddNewUser(c tele.Context) error {
	id := c.Sender().ID
	name := c.Sender().Username
	ID := int(id)
	err := s.Service.SaveInDB(name, ID)
	if err != nil {
		return err
	}
	err = c.Send("Отлично!")
	if err != nil {
		return err
	}
	return nil
}

var lastMessage *tele.Message

func (s *Handler) AddPlayer(c tele.Context) error {

	userID := c.Sender().ID
	chatID := c.Chat().ID
	err := s.Service.Inc(int(chatID), int(userID))
	if err != nil {
		s.Bot.Send(c.Chat(), "Игра уже идет")
		return nil
	}
	reply := "У вас есть 1 минута чтобы другие участники смогли присоединиться!"
	m := c.Message()
	exists := false
	if lastMessage != nil && lastMessage.Text == m.Text {
		exists = true
	}
	if exists {
		return nil
	}
	lastMessage = c.Message()
	s.Bot.Send(c.Chat(), reply)
	duration := 10 * time.Second
	timer := time.NewTimer(duration)

	go func() {
		<-timer.C
		reply := &tele.ReplyMarkup{}
		btn := reply.Data("Ready", "ready")
		reply.Inline(
			reply.Row(btn))
		s.Bot.Send(c.Chat(), "Нажми кнопку 'Ready'", reply)
	}()
	return nil
}

func (s *Handler) GiveCards(c tele.Context) error {
	userID := c.Sender().ID
	chatID := c.Chat().ID
	m, err := s.Service.AddInMap(int(chatID), int(userID))
	if err != nil {
		s.Bot.Send(c.Chat(), "Игра уже идет")
		return nil
	}
	for k, v := range m {
		if k == int(chatID) {
			for _, i := range v {
				if i.ID == int(userID) {
					for p, d := range i.Img {
						open, err := os.Open("/Users/dmitrydenisov/GolandProjects/Imaginarium/src/" + d.FileLocal)
						photo := &tele.Photo{File: tele.FromDisk(open.Name())}
						if err != nil {
							return nil
						}
						_, err = s.Bot.Send(c.Sender(), photo)
						if err != nil {
							return nil
						}
						btn := tele.InlineButton{
							Unique: strconv.Itoa(p),
							Text:   fmt.Sprint("Для ассоциации №" + strconv.Itoa(p)),
						}

						inlineKeys := [][]tele.InlineButton{
							[]tele.InlineButton{btn},
						}

						s.Bot.Send(c.Sender(), "Выберите фото:", &tele.ReplyMarkup{
							InlineKeyboard: inlineKeys,
						})
					}
				}
			}
		}
	}
	ready := s.Service.MapIsFull(int(chatID), int(userID))
	if ready {
		err = s.StartGame(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Handler) Association(c tele.Context) error {
	data := c.Message().Text
	if strings.HasPrefix(data, "/") {
		str, chat, err := s.Service.Association(data, int(c.Sender().ID))
		if err != nil {
			return err
		}
		chatID := tele.ChatID(chat)
		result := "Ассоциация была такая: " + str
		s.Bot.Send(chatID, result)
	}
	return nil
}

func (s *Handler) StartGame(c tele.Context) error {
	user, err := s.Service.StartG(int(c.Chat().ID))
	if err != nil {
		return err
	}
	res := "@" + user
	s.Bot.Send(c.Chat(), res)
	return nil
}

func (s *Handler) PhotoTake(c tele.Context) error {
	photoNumber := c.Data()
	userID := c.Sender().ID
	number, _ := strconv.Atoi(photoNumber[2:])
	chat, resPhoto, err := s.Service.TakePhoto(int(userID), number)
	if err != nil {
		return err
	}
	for _, ph := range resPhoto {
		for i, c := range ph.Img {
			open, err := os.Open("/Users/dmitrydenisov/GolandProjects/Imaginarium/src/" + c.FileLocal)
			if err != nil {
				return err
			}
			defer open.Close()
			phot := &tele.Photo{File: tele.FromDisk(open.Name())}
			chatID := tele.ChatID(chat)
			s.Bot.Send(chatID, phot)
			btn := tele.InlineButton{
				Unique: strconv.Itoa(i),
				Text:   fmt.Sprint("Голосование №" + strconv.Itoa(i)),
			}

			inlineKeys := [][]tele.InlineButton{
				[]tele.InlineButton{btn},
			}

			s.Bot.Send(chatID, "Выберите фото:", &tele.ReplyMarkup{
				InlineKeyboard: inlineKeys,
			})
		}
	}
	return nil

}
