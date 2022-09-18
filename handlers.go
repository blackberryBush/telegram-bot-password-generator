package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"strconv"
)

const (
	messageCommand = iota
	messageSticker
	messageChangeLength
	messageText
	messageUndefined = -1
)

func (b *Bot) getMessageType(message *tgbotapi.Message) int {
	if message.IsCommand() {
		return messageCommand
	}
	if message.Text == "" && message.Sticker != nil {
		return messageSticker
	}
	if b.changePassLength && message.Text != "" {
		return messageChangeLength
	}
	if message.Text != "" {
		return messageText
	}
	return messageUndefined
}

func (b *Bot) handleCommand(message *tgbotapi.Message, chatID int) error {
	keyboardDefault := tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(`/generate`),
		tgbotapi.NewKeyboardButton(`/options`),
		tgbotapi.NewKeyboardButton(`/stop`))

	defaultMessage := "Howdy!"
	switch message.Command() {

	case "start":
		msg := tgbotapi.NewMessage(int64(chatID), defaultMessage)
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(keyboardDefault)
		go b.Pull(chatID, msg)

	case "stop":
		msg := tgbotapi.NewMessage(int64(chatID), `SeeYa!`)
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
		go b.Pull(chatID, msg)

	case "options":
		replyText := "Выберите настройки, затем пропишите /generate"
		b.showOptions(int64(chatID), replyText)

	case "generate":
		newPass := ""
		newPass, err := GeneratePassword(b.param.length, b.param.number, b.param.upperCase, b.param.lowerCase, b.param.specialCase)
		if err != nil {
			newPass = "Password could not be generated"
			return err
		}
		msg := tgbotapi.NewMessage(int64(chatID), newPass)
		msg.ReplyToMessageID = message.MessageID // отвечаемое сообщение
		go b.Pull(chatID, msg)

	default:
		return b.handleUnknown()
	}
	return nil
}

func (b *Bot) handleSticker(message *tgbotapi.Message, chatID int) error {
	names := []string{"C:/Users/Konst/Downloads/32938fee-1632-4004-a834-5b310a42ba80.webp"}
	////////////////////////////////////////////////////////////////////////////////////////
	photoBytes, err := os.ReadFile(names[0])
	if err != nil {
		return err
	}
	photoFileBytes := tgbotapi.FileBytes{
		Name:  "picture",
		Bytes: photoBytes,
	}
	go b.Pull(chatID, tgbotapi.NewSticker(message.Chat.ID, photoFileBytes))
	return err
}

func (b *Bot) handleChangeLength(message *tgbotapi.Message, chatID int) error {
	tempPassLen, err := strconv.Atoi(message.Text)
	replyText := "Настройки изменены. Если вы завершили настройку, просто пропишите /generate"
	if err == nil && tempPassLen > 0 && tempPassLen <= 500 {
		b.param.length = tempPassLen
	} else {
		replyText = "Произошла ошибка. Попробуйте ещё раз."
	}
	b.changePassLength = false
	if err = InsertData("UsersDB", "Options", chatID, b.param, b.changePassLength); err != nil {
		return err
	}
	b.showOptions(int64(chatID), replyText)
	return err
}

func (b *Bot) handleText(message *tgbotapi.Message, chatID int) {
	msg := tgbotapi.NewMessage(int64(chatID), message.Text)
	msg.ReplyToMessageID = message.MessageID // Reply
	go b.Pull(chatID, msg)
}

func (b *Bot) sendSimpleMessage(chatID int, message string) {
	msg := tgbotapi.NewMessage(int64(chatID), message)
	go b.Pull(chatID, msg)
}

func (b *Bot) handleMessage(message *tgbotapi.Message, chatID int) error {
	var err error
	switch b.getMessageType(message) {
	case messageCommand:
		err = b.handleCommand(message, chatID)
	case messageSticker:
		err = b.handleSticker(message, chatID)
	case messageChangeLength:
		err = b.handleChangeLength(message, chatID)
	case messageText:
		b.handleText(message, chatID)
	default:
		err = b.handleUnknown()
	}
	return err
}

// Respond to the callback query (inverse passNumber)
func (b *Bot) handlePassNumber(callback *tgbotapi.CallbackQuery, chatID int) {
	b.param.number = !b.param.number
	replyText := ""
	if b.param.number {
		replyText = "В пароле будут присутствовать цифры"
	} else {
		replyText = "В пароле больше НЕ будут присутствовать цифры"
	}
	reply := tgbotapi.NewCallback(callback.ID, replyText)
	go b.Pull(chatID, reply)
}

func (b *Bot) handlePassUpperCase(callback *tgbotapi.CallbackQuery, chatID int) {
	b.param.upperCase = !b.param.upperCase
	replyText := ""

	if b.param.upperCase {
		replyText = "В пароле будут присутствовать буквы верхнего регистра"
	} else {
		replyText = "В пароле больше НЕ будут присутствовать буквы верхнего регистра"
	}

	reply := tgbotapi.NewCallback(callback.ID, replyText)
	go b.Pull(chatID, reply)
}

func (b *Bot) handlePassLowerCase(callback *tgbotapi.CallbackQuery, chatID int) {
	b.param.lowerCase = !b.param.lowerCase
	replyText := ""
	if b.param.lowerCase {
		replyText = "В пароле будут присутствовать буквы нижнего регистра"
	} else {
		replyText = "В пароле больше НЕ будут присутствовать буквы нижнего регистра"
	}
	reply := tgbotapi.NewCallback(callback.ID, replyText)
	go b.Pull(chatID, reply)
}

func (b *Bot) handlePassSpecialCase(callback *tgbotapi.CallbackQuery, chatID int) {
	b.param.specialCase = !b.param.specialCase
	replyText := ""
	if b.param.specialCase {
		replyText = "В пароле будут присутствовать спец.символы"
	} else {
		replyText = "В пароле больше НЕ будут присутствовать спец.символы"
	}
	reply := tgbotapi.NewCallback(callback.ID, replyText)
	go b.Pull(chatID, reply)
}

func (b *Bot) handlePassLength(callback *tgbotapi.CallbackQuery, chatID int) {
	b.changePassLength = true
	replyText := "Введите длину пароля:"
	msg := tgbotapi.NewMessage(callback.From.ID, replyText)
	go b.Pull(chatID, msg)
}

func (b *Bot) handleUnknown() error {
	return fmt.Errorf("unknown item was received")
}

func (b *Bot) handleCallbackQuery(callback *tgbotapi.CallbackQuery, chatID int) error {
	switch callback.Data {
	case "passNumber":
		b.handlePassNumber(callback, chatID)
	case "passUpperCase":
		b.handlePassUpperCase(callback, chatID)
	case "passLowerCase":
		b.handlePassLowerCase(callback, chatID)
	case "passSpecialCase":
		b.handlePassSpecialCase(callback, chatID)
	case "passLength":
		b.handlePassLength(callback, chatID)
	default:
		return b.handleUnknown()
	}
	msg := tgbotapi.NewDeleteMessage(int64(chatID), callback.Message.MessageID)
	go b.Pull(chatID, msg)
	if err := InsertData("UsersDB", "Options", chatID, b.param, b.changePassLength); err != nil {
		return err
	}
	if b.changePassLength {
		return nil
	}
	replyText := "Настройки изменены. Если вы завершили настройку, просто пропишите /generate"
	b.showOptions(int64(chatID), replyText)
	return nil
}
