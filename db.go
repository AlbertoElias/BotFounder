package main

import (
	"fmt"
	"os"
	"strings"

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
	env := os.Environ()
	host := "localhost"
	for _, s := range env {
		arr := strings.Split(s, "=")
		if arr[0] == "POSTGRES_PORT_5432_TCP_ADDR" {
			host = arr[1]
		}
	}

	db, err := gorm.Open("postgres", fmt.Sprintf("sslmode=disable host=%s", host))
	db.AutoMigrate(&User{}, &Conversation{}, &Bot{})

	return &DB{db}, err
}
