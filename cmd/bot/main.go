package main

import (
	"Imaginarium/config"
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

	b.Handle("/start", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	b.Start()
}
