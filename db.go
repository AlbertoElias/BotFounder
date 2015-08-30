package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// DB Internal database structure
type DB struct {
	db gorm.DB
}

// Database schema
type (
	// User in Database
	User struct {
		ID             int32 `gorm:"primary_key"`
		Name           string
		Conversation   Conversation
		ConversationID int32
		Bots           []Bot
	}

	// Conversation is each one of the people or groups a bot is talking to
	Conversation struct {
		ID                     int32 `gorm:"primary_key"`
		BotID                  int32 `sql:"index"`
		TelegramConversationID string
	}

	// Bots that have registered with the founder
	Bot struct {
		ID            int32 `gorm:"primary_key"`
		UserID        int32 `sql:"index"`
		TelegramToken string
		LastUpdate    int
		Conversations []Conversation
	}

	Config struct {
		ID         int32 `gorm:"primary_key"`
		LastUpdate int
	}
)

func (db DB) AllTheBots() (bots []Bot) {
	db.db.Find(&bots)
	return bots
}

func (db DB) GetLastUpdate() int {
	config := new(Config)
	db.db.FirstOrCreate(config)
	return config.LastUpdate
}

func (db DB) SetLastUpdate(last int) {
	config := new(Config)
	db.db.First(config)
	config.LastUpdate = last
	db.db.Save(config)
}

func (db DB) NewUser(id string) User {
	user := User{Conversation: Conversation{TelegramConversationID: id}}
	db.db.NewRecord(user)
	db.db.Create(&user)

	return user
}

func (db DB) GetBot(id string) *Bot {
	bot := new(Bot)
	db.db.Debug().Where("id = ?", id).First(bot)
	fmt.Println("got bot", bot)
	return bot
}

func (db DB) GetBotWithToken(token string) *Bot {
	bot := new(Bot)
	db.db.Where("telegram_token = ?", token).First(bot)
	return bot
}

func (db DB) SetLastUpdateForBot(id string, lastUpdate int) {
	bot := new(Bot)
	db.db.Where("id = ?", id).First(bot)
	bot.LastUpdate = lastUpdate
	db.db.Save(bot)

}

func (db DB) CreateConversationForBot(botid string, convid string) *Conversation {
	b, _ := strconv.Atoi(botid)

	conv := new(Conversation)
	State.DB.db.FirstOrCreate(conv, Conversation{BotID: int32(b), TelegramConversationID: convid})
	return conv
}

func (db DB) GetConversation(id string) *Conversation {
	con := new(Conversation)
	db.db.Where("id = ?", id).First(con)
	return con
}

func (db DB) GetConversationWithTelegram(id string) *Conversation {
	con := new(Conversation)
	db.db.Where("telegram_conversation_id = ?", id).First(con)
	return con
}

/*
func (db DB) SetBot(token string) {
	fmt.Println(token)
}

func (db DB) SetConversation(ConvId string) {
	fmt.Println(ConvId)
}

func (db DB) SetUser(UserId string) {
	fmt.Println(UserId)
}*/

// SetupDb Connect with postgres
func SetupDb() (*DB, error) {

	host := os.Getenv("POSTGRES_PORT_5432_TCP_ADDR")

	db, err := gorm.Open("postgres", fmt.Sprintf("sslmode=disable host=%s", host))
	db.AutoMigrate(&User{}, &Conversation{}, &Bot{}, &Config{})

	return &DB{db}, err
}
