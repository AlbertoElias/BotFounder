package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// DB Internal database structure
type DB struct {
	db gorm.DB
}

// User in Database
type User struct {
	ID           string
	Conversation Conversation
	Bots         []Bot
}

// Conversation is each one of the people or groups a bot is talking to
type Conversation struct {
	ID                     string
	TelegramConversationID string
	PostID                 string
}

// Bot s can talk to
type Bot struct {
	ID            string
	TelegramToken string
	Conversation  []Conversation
	PostID        string //TODO: think about whether a bot can have different post ids and send different messages to different converations
}

// SetupDb Connect with postgres
func SetupDb() (*DB, error) {
	db, err := gorm.Open("postgres", "sslmode=disable")
	db.AutoMigrate(&User{}, &Conversation{}, &Bot{})

	return &DB{db}, err
}
