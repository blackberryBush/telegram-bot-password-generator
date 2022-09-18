package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"time"
)

// Тестами не покрывал
// В обработку ошибок не вдавался, всё просто летит в лог

func token() string {
	return "5764246376:AAFQNHlXzHbLt8q79N4NDtFrmk_dM-76gzM"
}

func (b *Bot) showOptions(chatID int64, replyText string) {
	keyboardSetParam := getKeyboardParam(b.param.length, b.param.number, b.param.upperCase, b.param.lowerCase, b.param.specialCase)
	msg := tgbotapi.NewMessage(chatID, replyText)
	msg.ReplyMarkup = keyboardSetParam
	go b.Pull(int(chatID), msg)
}

func getKeyboardParam(length int, number, upperCase, lowerCase, specialCase bool) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Цифры "+checker(number), "passNumber"),
			tgbotapi.NewInlineKeyboardButtonData("Заглавные буквы "+checker(upperCase), "passUpperCase")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Строчные буквы "+checker(lowerCase), "passLowerCase"),
			tgbotapi.NewInlineKeyboardButtonData("Спец.символы "+checker(specialCase), "passSpecialCase")),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Длина пароля: "+strconv.Itoa(length), "passLength")))
}

func checkLastSend(chatID int, messageTimes map[int]time.Time) bool {
	if val, ok := messageTimes[chatID]; ok {
		dt := time.Now()
		return dt.After(val.Add(time.Second))
	}
	return true
}

func (b *Bot) TimeStart() {
	messageTimes := make(map[int]time.Time)
	timer := time.NewTicker(time.Second / 30)
	for range timer.C {
		b.toSend.Range(func(i int, v ItemToSend) bool {
			if v.queue > 0 && checkLastSend(i, messageTimes) {
				err := b.Send(i)
				if err != nil {
					log.Println(err)
				}
				messageTimes[i] = time.Now()
				return false
			}
			return true
		})
	}
}

func (b *Bot) Run() {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.bot.GetUpdatesChan(u)
	for update := range updates {
		var err error
		chatID := update.SentFrom().ID
		PrintReceive(&update)
		if b.param, b.changePassLength, err = GetData("UsersDB", "Options", int(chatID)); err != nil {
			if err := InsertData("UsersDB", "Options", int(chatID), b.param, b.changePassLength); err != nil {
				log.Println(err)
			}
		}
		if update.Message != nil {
			err = b.handleMessage(update.Message, int(chatID))
		} else if update.CallbackQuery != nil {
			err = b.handleCallbackQuery(update.CallbackQuery, int(chatID))
		}
		if err != nil {
			log.Println(err)
		}
	}

}

func main() {
	botAPI, err := tgbotapi.NewBotAPI(token())
	if err != nil {
		log.Fatal(err)
	}
	b := NewBot(botAPI)
	// Start timer to send messages&callbacks
	go b.TimeStart()
	// Create DB and table if not exist
	CreateDatabase("UsersDB")
	CreateTable("UsersDB", "Options")
	// Start checking for updates
	b.Run()
}
