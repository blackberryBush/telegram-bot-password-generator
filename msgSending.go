package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) Send(chatID int) (err error) {
	b.toSend.QueueDec(chatID)
	item, ok := b.toSend.Load(chatID)
	if !ok {
		return fmt.Errorf("reading error")
	}
	data := <-item.data
	if item.queue > 7 {
		return nil
	}
	if item.queue > 5 {
		data = tgbotapi.NewMessage(int64(chatID), "Не флуди!")
	}
	PrintSent(&data)
	switch data.(type) {
	case tgbotapi.MessageConfig:
		_, err = b.bot.Send(data)
	case tgbotapi.CallbackConfig, tgbotapi.DeleteMessageConfig:
		_, err = b.bot.Request(data)
	default:
		err = fmt.Errorf("undefined type")
	}
	return err
}

func (b *Bot) Pull(chatID int, c tgbotapi.Chattable) {
	if v, ok := b.toSend.Load(chatID); ok {
		b.toSend.QueueInc(chatID)
		v.data <- c
		b.toSend.StoreData(chatID, v.data)
	} else {
		item := *NewItemToSend()
		item.data <- c
		item.queue++
		b.toSend.Store(chatID, item)
	}
}
