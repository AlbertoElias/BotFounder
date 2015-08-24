package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Server internal structure
type Server struct {
	Server *gin.Engine
}

// SetupServer configures and returns a server instance
func SetupServer() *Server {
	router := gin.Default()

	router.GET("/post/:id/:message/:conversation", func(c *gin.Context) {
		botId := c.Param("id")
		message := c.Param("message")
		conversation := c.Param("conversation")

		bot := State.DB.GetBot(botId).Bot()
		bot = State.Bot
		success := <-bot.SendMessage(message, conversation) // Function returns a channel, and we wait for the channel to send something
		c.String(http.StatusOK, "Hello %s. Message is: %s", bot.Token, success)
	})

	router.GET("/post/:id/:message", func(c *gin.Context) {
		botId := c.Param("id")
		message := c.Param("message")

		bot := State.DB.GetBot(botId).Bot()
		success := <-bot.SendMessage(message, "") // Function returns a channel, and we wait for the channel to send something
		c.String(http.StatusOK, "Hello %s. Message is: %s", bot.Token, success)
	})

	// Webhook url, gets "update" and passes it to Bot.parseUpdate()
	router.GET(fmt.Sprintf("/%s", os.Getenv("FOUNDERBOT_TOKEN")), func(c *gin.Context) {

	})

	return &Server{router}
}

// Run starts the server
func (s *Server) Run() {
	s.Server.Run(":3000")
}
