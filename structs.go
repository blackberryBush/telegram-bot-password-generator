package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PasswordParam struct {
	length                                    int
	number, upperCase, lowerCase, specialCase bool
}

type Bot struct {
	bot              *tgbotapi.BotAPI
	param            *PasswordParam
	changePassLength bool
	toSend           *ItemsToSend
}

type ItemToSend struct {
	queue int
	data  chan tgbotapi.Chattable
}

func NewItemToSend() *ItemToSend {
	return &ItemToSend{
		queue: 0,
		data:  make(chan tgbotapi.Chattable, 1),
	}
}

func NewPasswordParam(length int, number, upperCase, lowerCase, specialCase bool) *PasswordParam {
	return &PasswordParam{
		length:      length,
		number:      number,
		upperCase:   upperCase,
		lowerCase:   lowerCase,
		specialCase: specialCase,
	}
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {
	return &Bot{
		bot:              bot,
		param:            NewPasswordParam(20, true, true, true, false),
		changePassLength: false,
		toSend:           NewItemsToSend(),
	}
}
