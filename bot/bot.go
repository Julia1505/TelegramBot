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

	WebhookURL = "https://1a34-188-255-34-137.eu.ngrok.io"
)

type TelegramBot struct {
	bot *tgbotapi.BotAPI
}

func startTaskBot(ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatalf("NewBotApi failed: %s", err)
		return err
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

	bot.Debug = true
	updates := bot.ListenForWebhook("/")

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

	t := &TelegramBot{
		bot: bot,
	}

	for {
		select {
		case update := <-updates:
			log.Printf("upd: %#v\n", update)

			t.HandleUpdates(storage, update)
		case <-ctx.Done():
			return nil
		}

	}

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	err := startTaskBot(ctx)
	if err != nil {
		//nolint:govet
		fmt.Printf("startTaskBot error: %s", err)
	}

	defer cancel()

}
