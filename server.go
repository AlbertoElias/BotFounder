package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server internal structure
type Server struct {
	Server *gin.Engine
}

// SetupServer configures and returns a server instance
func SetupServer() *Server {
	router := gin.Default()

	router.GET("/bot/:botid/:conversation", func(c *gin.Context) {
		botId := c.Param("botid")
		message := c.Query("message")
		conversation := c.Param("conversation")

		bot := State.DB.GetBot(botId).Bot()
		success := <-bot.SendMessage(message, conversation) // Function returns a channel, and we wait for the channel to send something
		c.String(http.StatusOK, "Hello %s. Message is: %s", bot.Token, success)
	})

	router.GET("/bot/:botid", func(c *gin.Context) {
		botId := c.Param("botid")
		message := c.Query("message")

		bot := State.DB.GetBot(botId).Bot()
		success := <-bot.SendMessage(message, "all") // Function returns a channel, and we wait for the channel to send something
		c.String(http.StatusOK, "Hello %s. Message is: %s", bot.Token, success)
	})

	router.GET("/s/:convid", func(c *gin.Context) {
		convid := c.Param("convid")
		message := c.Query("message")
		conver := State.DB.GetConversation(convid)

		if conver == nil {
			c.String(http.StatusNotFound, "Not found conversation")
		} else {
			success := <-State.Bot.SendMessage(message, conver.TelegramConversationID)
			c.String(http.StatusOK, "sent message %s %s", conver.TelegramConversationID, success)
		}
	})

	/*router.GET("/")

	// Webhook url, gets "update" and passes it to Bot.parseUpdate()
	router.GET(fmt.Sprintf("/%s", os.Getenv("FOUNDERBOT_TOKEN")), func(c *gin.Context) {

	})*/

	return &Server{router}
}

// Run starts the server
func (s *Server) Run() {
	s.Server.Run(":3000")
}
