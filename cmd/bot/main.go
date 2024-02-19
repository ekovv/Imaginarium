package main

import (
	"Imaginarium/config"
	"Imaginarium/internal/handler"
	"Imaginarium/internal/service"
	"Imaginarium/internal/storage"
	"gopkg.in/telebot.v3"
	"log"
)

func main() {
	conf := config.New()
	st, err := storage.New(conf)
	if err != nil {
		log.Fatal("bad")
	}
	sr := service.NewService(st)
	h, err := handler.NewHandler(sr, conf)
	if err != nil {
		return
	}
	h.Bot.Handle("/login", h.AddNewUser)
	h.Bot.Handle("/start", h.Start)
	h.Bot.Handle(telebot.OnText, h.Association)
	h.Bot.Handle(telebot.OnCallback, h.HandleButton)
	h.Bot.Start()
}
