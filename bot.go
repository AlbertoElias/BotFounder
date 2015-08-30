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
func NewBot(token string, id string, lastUpdate int) *TelegramBot {
	return &TelegramBot{
		Token:      token,
		DatabaseID: id,
		LastUpdate: lastUpdate,
	}
}

// Bot creates a bot from Bot database object
func (b Bot) Bot() *TelegramBot {
	return NewBot(b.TelegramToken, fmt.Sprintf("%d", b.ID), b.LastUpdate)
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
func (b *TelegramBot) PollConversationsEvery(seconds time.Duration) {

	fmt.Println("Starting updates for bot", b.DatabaseID)
	ticker := time.NewTicker(time.Second * seconds)
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
				if b.DatabaseID == "founder" {
					b.FounderUpdates(messages)
				} else {
					b.Updates(messages)
				}
			}
		}
	}()
}

func (b *TelegramBot) Updates(messages []Message) {
	for _, m := range messages {
		fmt.Println("update", m)
		converId := fmt.Sprintf("%d", m.Chat.Id)
		conver := State.DB.CreateConversationForBot(b.DatabaseID, converId) // Done

		// Send confirmation message to owner
		user := new(User)
		userconv := new(Conversation)
		bot := State.DB.GetBot(b.DatabaseID)
		fmt.Println("bot", bot)
		State.DB.db.Model(bot).Related(user)
		fmt.Println("the user", user)
		State.DB.db.Model(user).Related(userconv)
		fmt.Println("now conv", userconv)
		State.Bot.SendMessage(fmt.Sprintf("New conversation. You can send specific messages using this URL: %d", conver.ID), userconv.TelegramConversationID)
	}
}

func (b *TelegramBot) FounderUpdates(messages []Message) {
	for _, m := range messages {

		if strings.Contains(m.Text, "start") {
			converId := fmt.Sprintf("%d", m.Chat.Id)
			conver := State.DB.GetConversationWithTelegram(converId)
			fmt.Println(conver)
			if *conver == *new(Conversation) {

				user := State.DB.NewUser(converId)
				conver = &user.Conversation
				fmt.Println(conver)
			}

			b.SendMessage(fmt.Sprintf("Hey, welcome! You can send messages with this URL: http://localhost:3000/s/%d", conver.ID), converId)
		} else if strings.Contains(m.Text, "token") {
			converId := fmt.Sprintf("%d", m.Chat.Id)
			conver := State.DB.GetConversationWithTelegram(converId)
			fmt.Println(conver)
			if *conver == *new(Conversation) {
				b.SendMessage("Hey, you need to /start before registering your bot.", converId)
			} else {
				user := new(User)
				State.DB.db.Model(conver).Related(user)
				strs := strings.Split(m.Text, " ")
				if len(strs) > 1 {
					token := strs[1]
					bot := Bot{UserID: user.ID, TelegramToken: token}
					State.DB.db.FirstOrCreate(&bot, bot)
					b.SendMessage(fmt.Sprintf("Bot registered! You can make the bot send messages with this URL: http://localhost:3000/bot/%d", bot.ID), converId)
					bot.Bot().PollConversationsEvery(60)
				}
			}
		}
	}
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

	for _, u := range tResponse.Result {
		if u.Id > b.LastUpdate {
			messages = append(messages, u.Message)
			b.LastUpdate = u.Id
		}
	}
	b.SetLastUpdate(b.LastUpdate)

	return messages
}

func (b *TelegramBot) SetLastUpdate(update int) {
	if b.DatabaseID == "founder" {
		State.DB.SetLastUpdate(b.LastUpdate)
	} else {
		State.DB.SetLastUpdateForBot(b.DatabaseID, update)
	}
}
