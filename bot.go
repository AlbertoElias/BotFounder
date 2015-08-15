package main

import (
	"fmt"
	"net/http"
	"time"
)

// TelegramBot talks to telegram bot API
type TelegramBot struct {
	Token string
}

// NewBot creates telegram bot object and makes it ready to talk to the API
func NewBot(token string) *TelegramBot {
	return &TelegramBot{token}
}

// Bot creates a bot from Bot database object
func (b *Bot) Bot() *TelegramBot {
	return NewBot(b.TelegramToken)
}

func (b *TelegramBot) request(method string, path string) {
	url := fmt.Sprintf("https://api.telegram.org/bot/%s/%s", b.Token, path)
	http.Get(url)
}

// SendMessage to the wanted conversation
func (b *TelegramBot) SendMessage(text string, conversation string) <-chan bool {
	ch := make(chan bool)

	go func() {
		<-time.Tick(5 * time.Second)
		ch <- true
	}()

	return ch
}

// SetupWebhook to receive updates from telegram
func (b *TelegramBot) SetupWebhook() {

}
