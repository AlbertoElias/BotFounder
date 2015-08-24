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
	}

	// Bots that have registered with the founder
	Bot struct {
		ID            string
		TelegramToken string
		Conversations []Conversation
	}
)

func (db DB) GetBot(id string) *Bot {
	return &Bot{
		ID:            id,
		TelegramToken: "1232533453sfdgfdg",
	}
}

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
