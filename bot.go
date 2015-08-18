package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
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
func (b Bot) Bot() *TelegramBot {
	return NewBot(b.TelegramToken)
}

func (b TelegramBot) request(method string, params map[string]string) {
	var paramsString bytes.Buffer
	for key, value := range params {
		paramsString.WriteString(fmt.Sprintf("%s=%s&", key, value))
	}
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s?%s", b.Token, method, paramsString.String())
	http.Get(url)
}

// SendMessage to the wanted conversation
func (b TelegramBot) SendMessage(text string, conversation string) <-chan bool {
	ch := make(chan bool)

	go func() {
		<-time.Tick(5 * time.Second)
		ch <- true
	}()

	return ch
}

// SetupWebhook to receive updates from telegram. We need an SSL cert though
func (b TelegramBot) SetupWebhook() {
	b.request("setWebhook", map[string]string{
		"url": fmt.Sprintf("http://oururl.com/%s", os.Getenv("FOUNDERBOT_TOKEN")),
	})
}

// Longpolling of the "getUpdates" method
func (b TelegramBot) pollConversations() {

}
