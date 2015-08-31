package main

import (
	"fmt"
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
	router.LoadHTMLGlob("templates/*")

	router.GET("/bot/:botid/:conversation", func(c *gin.Context) {
		botId := c.Param("botid")
		message := c.Query("message")
		conversation := c.Param("conversation")

		bot := State.DB.GetBot(botId).Bot()
		success := <-bot.SendMessage(message, conversation) // Function returns a channel, and we wait for the channel to send something
		c.String(http.StatusOK, "Hello %s. Message is: %s", bot.Token, success)
	})

	router.GET("/bot/:botid", func(c *gin.Context) {
		botid := c.Param("botid")
		bot := State.DB.GetBot(botid)
		if bot.ID == 0 {
			c.String(http.StatusNotFound, "Not found bot")
		} else {
			c.HTML(http.StatusOK, "sendMessage.tmpl", gin.H{"ok": c.Query("ok"), "postTo": fmt.Sprintf("/bot/%s", botid)})
		}
	})

	router.POST("/bot/:botid", func(c *gin.Context) {
		botId := c.Param("botid")
		message := c.PostForm("message")

		bot := State.DB.GetBot(botId)
		convs := []Conversation{}
		State.DB.db.Model(bot).Related(&convs)
		for _, c := range convs {
			bot.Bot().SendMessage(message, c.TelegramConversationID)
		}
		c.String(http.StatusOK, "ok")
		//c.Redirect(http.StatusFound, fmt.Sprintf("/bot/%s?ok=sending...", botId))
	})

	router.GET("/s/:convid", func(c *gin.Context) {
		convid := c.Param("convid")
		conver := State.DB.GetConversation(convid)
		if *conver == *new(Conversation) {
			c.String(http.StatusNotFound, "Not found conversation")
		} else {
			c.HTML(http.StatusOK, "sendMessage.tmpl", gin.H{"ok": c.Query("ok"), "postTo": fmt.Sprintf("/s/%s", convid)})
		}
	})

	router.POST("/s/:convid", func(c *gin.Context) {
		convid := c.Param("convid")
		message := c.PostForm("message")
		conver := State.DB.GetConversation(convid)

		if *conver == *new(Conversation) {
			c.String(http.StatusNotFound, "Not found conversation")
		} else {
			<-State.Bot.SendMessage(message, conver.TelegramConversationID)
			//c.Redirect(http.StatusFound, fmt.Sprintf("/s/%s?ok=sent", convid))
			c.String(http.StatusOK, "ok")
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
