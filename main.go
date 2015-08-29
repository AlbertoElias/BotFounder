package main

import (
	"fmt"
	"os"
)

// State takes care of storing global dependencies of the project
var State = struct {
	*DB
	*Server
	Bot *TelegramBot
}{}

func main() {
	db, err := SetupDb()
	panicOnErr(err)
	State.DB = db
	defer db.db.Close()

	State.Bot = NewBot(os.Getenv("FOUNDERBOT_TOKEN"), true)
	State.Bot.pollConversations()
	play()
	State.Server = SetupServer()
	State.Server.Run()
}

func play() {
	user := State.DB.NewUser("860578")
	fmt.Println(user)
}

// HandleError decides what to do with an error. Right now it just panics.
func HandleError(err error) {
	panicOnErr(err)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
