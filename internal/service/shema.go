package service

import tele "gopkg.in/telebot.v3"

type Gamers struct {
	ID  int
	Img []*tele.Photo
}

type Voting struct {
	IDWin        int
	NicknameVote string
	NicknameWin  string
	Count        int
}
