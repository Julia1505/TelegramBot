package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

func (t *TelegramBot) HandleUpdates(st *TaskStorage, update tgbotapi.Update) {
	if message := update.Message; message.Text != "" {
		parseMessage := strings.Fields(message.Text)

		var respMessage string
		switch parseMessage[0] {
		case "/tasks":
			respMessage = st.Show(update, User{}, User{})
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, respMessage))
		case "/new":
			respMessage = st.AddTask(update)
			t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, respMessage))
			//case "/my":
			//
			//case "/owner":

		}

	}
}

func ShowAll(t *TelegramBot, update tgbotapi.Update) {

}
