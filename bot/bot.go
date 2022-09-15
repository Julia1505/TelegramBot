package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
)

var (
	BotToken = "5130893085:AAFxPRK32MNUj8A1pBbvTuJMN1kYLOc5ZkM"

	WebhookURL = "https://c445-188-255-34-137.eu.ngrok.io"
)

type TelegramBot struct {
	bot *tgbotapi.BotAPI
}

func startTaskBot(ctx context.Context) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatalf("NewBotApi failed: %s", err)
		return nil, err
	}

	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	wh, err := tgbotapi.NewWebhook(WebhookURL)
	if err != nil {
		log.Fatalf("NewWebhook failed: %s", err)
	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Fatalf("SetWebhook failed: %s", err)
	}

	return &TelegramBot{
			bot: bot,
		},
		nil
}

func (t *TelegramBot) Run(debug bool) error {
	t.bot.Debug = debug
	updates := t.bot.ListenForWebhook("/")

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all is working"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go func() {
		log.Fatalln("http err:", http.ListenAndServe(":"+port, nil))
	}()
	fmt.Println("start listen :" + port)

	storage := CreateTaskCollection()

	for update := range updates {
		log.Printf("upd: %#v\n", update)

		t.HandleUpdates(storage, update)
		//mes := update.Message.Text
		//t.bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Купи слона "+mes+update.Message.From.UserName))
	}

	return nil
}

func main() {
	bot, err := startTaskBot(context.Background())
	bot.Run(true)

	if err != nil {
		panic(err)
	}
}
