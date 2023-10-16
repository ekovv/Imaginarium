package handler

import (
	"Imaginarium/internal/service"
	tele "gopkg.in/telebot.v3"
)

type Handler struct {
	service service.Inter
}

func NewHandler(service service.Inter) *Handler {
	return &Handler{service: service}
}

func (s *Handler) AddNewUser(c tele.Context) error {
	err := s.service.SaveInDB(c.Data())
	if err != nil {
		return nil
	}
	c.Send("Отлично!")
	return nil
}
