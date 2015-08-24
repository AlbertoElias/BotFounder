package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// TelegramBot talks to telegram bot API
type TelegramBot struct {
	Token string
}

// Structure corresponding to the "update" response
type TelegramUpdate struct {
	first_name string
}

// NewBot creates telegram bot object and makes it ready to talk to the API
func NewBot(token string) *TelegramBot {
	return &TelegramBot{
		Token: token,
	}
}

// Bot creates a bot from Bot database object
func (b Bot) Bot() *TelegramBot {
	return NewBot(b.TelegramToken)
}

func (b TelegramBot) request(method string, params map[string]string) []byte {
	var paramsString bytes.Buffer
	for key, value := range params {
		paramsString.WriteString(fmt.Sprintf("%s=%s&", key, value))
	}
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s?%s", b.Token, method, paramsString.String())
	fmt.Println(url)

	resp, err := http.Get(url)
	panicOnErr(err)
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	return respBody
}

// SendMessage to the wanted conversation
func (b TelegramBot) SendMessage(text string, conversation string) <-chan []byte {
	ch := make(chan []byte)

	go func() {
		ch <- b.request("getMe", make(map[string]string))
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

// Parse telegram conversation updates
func parseUpdate(update []byte) *TelegramUpdate {
	var tUpdate TelegramUpdate
	err := json.Unmarshal(update, &tUpdate)
	panicOnErr(err)

	return &tUpdate
}
