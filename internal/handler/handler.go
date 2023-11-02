package handler

import (
	"Imaginarium/config"
	"Imaginarium/internal/service"
	tele "gopkg.in/telebot.v3"
	"os"
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
	}
	return nil
}

func (s *Handler) AddNewUser(c tele.Context) error {
	id := c.Sender().ID
	ID := int(id)
	err := s.Service.SaveInDB(ID)
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
	duration := 3 * time.Second
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
					for _, d := range i.Img {
						open, err := os.Open("/Users/dmitrydenisov/GolandProjects/Imaginarium/src/" + d.FileLocal)
						photo := &tele.Photo{File: tele.FromDisk(open.Name())}
						if err != nil {
							return nil
						}
						_, err = s.Bot.Send(c.Sender(), photo)
						if err != nil {
							return nil
						}
					}
				}
			}
		}
	}
	return nil
}

func (s *Handler) Association(c tele.Context) error {
	data := c.Message().Text
	if !strings.HasPrefix(data, "/") {
		s.Bot.Send(c.Sender(), "Напиши мне ассоциацию начиная с /")
	}
	res, err := s.Service.Association(data, int(c.Sender().ID))
	if err != nil {
		return err
	}
	result := "Ассоциация была такая: " + res
	s.Bot.Send(c.Chat(), result)

	return nil
}
