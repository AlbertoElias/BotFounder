package main

import "os"

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

	State.Bot = NewBot(os.Getenv("FOUNDERBOT_TOKEN"))
	State.Bot.pollConversations()

	State.Server = SetupServer()
	State.Server.Run()
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
