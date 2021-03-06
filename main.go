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
	defer db.db.Close()

	State.Bot = NewBot(os.Getenv("FOUNDERBOT_TOKEN"), "founder", State.DB.GetLastUpdate())
	State.Bot.PollConversationsEvery(1)

	for _, bot := range State.DB.GetAllBots() {
		bot.Bot().PollConversationsEvery(6)
	}

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
