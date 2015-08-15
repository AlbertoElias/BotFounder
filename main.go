package main

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

	State.Server = SetupServer()
	State.Bot = NewBot("whatever the token is")

	go State.Server.Run()
	select {} // Keep the "main thread" busy waiting for nothing so program does not exit
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
