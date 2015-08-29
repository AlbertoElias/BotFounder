package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// TelegramBot talks to telegram bot API
type TelegramBot struct {
	Token      string
	DatabaseID string
	LastUpdate int
}

// NewBot creates telegram bot object and makes it ready to talk to the API
func NewBot(token string, ourBot bool) *TelegramBot {
	bot := &TelegramBot{
		Token: token,
	}

	if ourBot {
		bot.DatabaseID = "founder"
		bot.LastUpdate = State.DB.GetLastUpdate()
	}

	return bot
}

// Bot creates a bot from Bot database object
func (b Bot) Bot() *TelegramBot {
	return NewBot(b.TelegramToken, false)
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
	panicOnErr(err)

	return respBody
}

// SendMessage to the wanted conversation
func (b TelegramBot) SendMessage(text string, conversation string) <-chan []byte {
	ch := make(chan []byte)

	conversations := []string{}
	if conversation == "all" {
		for _, conv := range State.DB.GetBot(b.DatabaseID).Conversations {
			conversations = append(conversations, conv.TelegramConversationID)
		}
	} else {
		conversations = []string{conversation}
	}

	for _, conv := range conversations {
		go func() {
			ch <- b.request("sendMessage", map[string]string{
				"chat_id": conv,
				"text":    text,
			})
		}()
	}
	return ch
}

// SetupWebhook to receive updates from telegram. We need an SSL cert though
func (b TelegramBot) SetupWebhook() {
	b.request("setWebhook", map[string]string{
		"url": fmt.Sprintf("http://oururl.com/%s", os.Getenv("FOUNDERBOT_TOKEN")),
	})
}

// Longpolling of the "getUpdates" method
func (b *TelegramBot) pollConversations() {
	ticker := time.NewTicker(time.Second * 3)
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
				messages := b.parseUpdate(update)
				for _, m := range messages {
					if strings.Contains(m.Text, "start") {

						b.SendMessage("hello my friend", fmt.Sprintf("%d", m.Chat.Id))
					}
				}
			}
		}
	}()
}

type UpdateResult struct {
	Id      int     `json:"update_id"`
	Message Message `json:"message"`
}

type Message struct {
	Id   int    `json:"message_id"`
	From Sender `json:"from"`
	Chat Sender `json:"chat"`
	Text string `json:"text"`
	Date int    `json:"date"`
}

type Sender struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

// Parse telegram conversation updates
func (b *TelegramBot) parseUpdate(response []byte) []Message {
	tResponse := struct {
		Ok     bool           `json:"ok"`
		Result []UpdateResult `json:"result"`
	}{}
	err := json.Unmarshal(response, &tResponse)
	panicOnErr(err)

	if !tResponse.Ok {
		return nil
	}

	messages := []Message{}

	fmt.Println(b.LastUpdate)
	for _, u := range tResponse.Result {
		if u.Id > b.LastUpdate {
			messages = append(messages, u.Message)
			b.LastUpdate = u.Id
		}
	}
	State.DB.SetLastUpdate(b.LastUpdate)

	return messages
}
