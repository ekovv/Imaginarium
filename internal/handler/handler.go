package handler

import (
	"Imaginarium/internal/service"
	tele "gopkg.in/telebot.v3"
	"log"
)

type Handler struct {
	service service.Inter
}

func NewHandler(service service.Inter) *Handler {
	return &Handler{service: service}
}

func (s *Handler) AddNewUser(c tele.Context) error {
	nickname := c.Sender().Username
	log.Println("Никнейм пользователя:", nickname)
	err := s.service.SaveInDB(nickname)
	if err != nil {
		return nil
	}
	c.Send("Отлично!")
	return nil
}
