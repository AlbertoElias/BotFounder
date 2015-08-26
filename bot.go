package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// TelegramBot talks to telegram bot API
type TelegramBot struct {
	Token string
}

// Structure corresponding to all responses
type TelegramResponse struct {
	ok     bool
	result []interface{}
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

	resp, err := http.Get(url)
	panicOnErr(err)
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	panicOnErr(err)

	return respBody
}

// SendMessage to the wanted conversation
func (b TelegramBot) SendMessage(text string, conversation string) <-chan []byte {
	ch := make(chan []byte)

	go func() {
		ch <- b.request("sendMessage", map[string]string{
			"chat_id": conversation,
			"text":    text,
		})
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
	ticker := time.NewTicker(time.Second * 30)
	updatesChannel := make(chan []byte)
	go func() {
		for _ = range ticker.C {
			updatesChannel <- b.request("getUpdates", make(map[string]string))
		}
	}()

	go func() {
		for {
			select {
			case update := <-updatesChannel:
				telegramUpdate := parseUpdate(update)
				fmt.Println(telegramUpdate)
			}
		}
	}()
}

// Parse telegram conversation updates
func parseUpdate(response []byte) *TelegramUpdate {
	var tResponse TelegramResponse
	err := json.Unmarshal(response, &tResponse)
	panicOnErr(err)

	fmt.Println(tResponse)
	fmt.Println(tResponse.ok)
	if tResponse.ok == true {
		var tUpdate TelegramUpdate
		fmt.Println(tResponse.result)
		// err := json.Unmarshal(update, &tUpdate)
		panicOnErr(err)

		return &tUpdate
	}

	return nil
}
