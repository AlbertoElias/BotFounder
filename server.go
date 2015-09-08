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

	router.GET("/bot/:botid/:convid", func(c *gin.Context) {
		botid := c.Param("botid")
		convid := c.Param("convid")
		bot := State.DB.GetBot(botid)

		if bot.ID == 0 {
			c.String(http.StatusNotFound, "Bot not found")
		} else {
			c.HTML(http.StatusOK, "sendMessage.tmpl", gin.H{"ok": c.Query("ok"), "postTo": fmt.Sprintf("/bot/%s/%s", botid, convid)})
		}
	})

	router.POST("/bot/:botid/:convid", func(c *gin.Context) {
		botId := c.Param("botid")
		convid := c.Param("convid")
		message := c.PostForm("message")
		conver := State.DB.GetConversation(convid)

		bot := State.DB.GetBot(botId)
		bot.Bot().SendMessage(message, conver.TelegramConversationID)
		c.String(http.StatusOK, "ok")
	})

	router.GET("/bot/:botid", func(c *gin.Context) {
		botid := c.Param("botid")
		bot := State.DB.GetBot(botid)
		if bot.ID == 0 {
			c.String(http.StatusNotFound, "Bot not found")
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
	})

	router.GET("/founderbot/:convid", func(c *gin.Context) {
		convid := c.Param("convid")
		conver := State.DB.GetConversation(convid)
		if *conver == *new(Conversation) {
			c.String(http.StatusNotFound, "Conversation not found")
		} else {
			c.HTML(http.StatusOK, "sendMessage.tmpl", gin.H{"ok": c.Query("ok"), "postTo": fmt.Sprintf("/founderbot/%s", convid)})
		}
	})

	router.POST("/founderbot/:convid", func(c *gin.Context) {
		convid := c.Param("convid")
		message := c.PostForm("message")
		conver := State.DB.GetConversation(convid)

		if *conver == *new(Conversation) {
			c.String(http.StatusNotFound, "Conversation not found")
		} else {
			<-State.Bot.SendMessage(message, conver.TelegramConversationID)
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
