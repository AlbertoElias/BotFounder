package main

import (
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
		ID           string
		Conversation Conversation
		Bots         []Bot
	}

	// Conversation is each one of the people or groups a bot is talking to
	Conversation struct {
		ID                     string
		TelegramConversationID string
		PostID                 string
	}

	// Bot s can talk to
	Bot struct {
		ID            string
		TelegramToken string
		Conversation  []Conversation
		PostID        string //TODO: think about whether a bot can have different post ids and send different messages to different converations
	}
)

// SetupDb Connect with postgres
func SetupDb() (*DB, error) {
	db, err := gorm.Open("postgres", "sslmode=disable")
	db.AutoMigrate(&User{}, &Conversation{}, &Bot{})

	return &DB{db}, err
}
