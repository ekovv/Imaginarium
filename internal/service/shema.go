package service

import tele "gopkg.in/telebot.v3"

type Gamers struct {
	ID  int
	Img []*tele.Photo
}

type Voting struct {
	ID       int
	Nickname string
	Count    int
}
