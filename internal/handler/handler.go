package handler

import (
	"Imaginarium/config"
	"Imaginarium/internal/service"
	tele "gopkg.in/telebot.v3"
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
func (s *Handler) AddPlayer(c tele.Context) error {
	id := c.Sender().ID
	ID := int(id)
	err := s.Service.AddInMap(ID)
	if err != nil {
		return err
	}

	_, err = s.Bot.Send(c.Sender(), "Ты зарегался в игру в чате")
	if err != nil {
		return err
	}
	_, err = s.Bot.Send(c.Chat(), "Ты авторизовался в игру")
	if err != nil {
		return err
	}
	return nil
}
