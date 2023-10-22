package handler

import (
	"Imaginarium/config"
	"Imaginarium/internal/service"
	tele "gopkg.in/telebot.v3"
	"os"
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
	newID := int(id)
	err := s.Service.Inc(newID)
	if err != nil {
		return err
	}
	reply := "У вас есть 1 минута чтобы другие участники смогли присоединиться!"
	s.Bot.Send(c.Chat(), reply)

	// Запускаем таймер на 5 секунд
	duration := 15 * time.Second
	timer := time.NewTimer(duration)
	// Горутина для обработки события истечения времени таймера
	go func() {
		<-timer.C // Ждем истечения таймера
		reply := &tele.ReplyMarkup{}

		btn := reply.Data("Ready", "button_callback")
		reply.Inline(
			reply.Row(btn),
		)
		s.Bot.Send(c.Chat(), "Нажми кнопку 'Ready'", reply)
	}()
	return nil
}

func (s *Handler) GiveCards(c tele.Context) error {
	m := s.Service.AddInMap()
	id := c.Sender().ID
	for k, v := range m {
		for _, i := range v {
			newK := int64(k)
			if id == newK {
				open, err := os.Open("/Users/dmitrydenisov/GolandProjects/Imaginarium/src/" + i.FileLocal)
				photo := &tele.Photo{File: tele.FromDisk(open.Name())}
				if err != nil {
					return err
				}
				_, err = s.Bot.Send(c.Sender(), photo)
				if err != nil {
					return err
				}
			}
		}

	}
	return nil
}
