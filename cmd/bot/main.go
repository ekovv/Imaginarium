package main

import (
	"Imaginarium/config"
	"Imaginarium/internal/handler"
	"Imaginarium/internal/service"
	"Imaginarium/internal/storage"
)

func main() {
	conf := config.New()
	st := storage.NewStorage(conf)
	sr := service.NewService(*st)
	h, err := handler.NewHandler(sr, conf)
	if err != nil {
		return
	}
	h.Bot.Handle("/login", h.AddNewUser)
	h.Bot.Handle("/play", h.AddPlayer)
	h.Bot.Start()
}
