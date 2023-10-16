package main

import (
	"Imaginarium/config"
	"Imaginarium/internal/handler"
	"Imaginarium/internal/service"
	"Imaginarium/internal/storage"
	tele "gopkg.in/telebot.v3"
	"log"
	"time"
)

func main() {
	conf := config.New()
	pref := tele.Settings{
		Token:  conf.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}
	st := storage.NewStorage(conf)
	sr := service.NewService(*st)
	h := handler.NewHandler(sr)

	b.Handle("/login", h.AddNewUser)
	b.Start()
}
