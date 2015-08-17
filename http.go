package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Server internal structure
type Server struct {
	Server *gin.Engine
}

// SetupServer configures and returns a server instance
func SetupServer() *Server {
	router := gin.Default()

	router.GET("/post/:id", func(c *gin.Context) {
		name := c.Param("id")

		success := <-State.Bot.SendMessage("lola", "12345") // Function returns a channel, and we wait for the channel to send something
		c.String(http.StatusOK, "Hello %s. It went %s", name, success)
	})

	return &Server{router}
}

// Run starts the server
func (s *Server) Run() {
	s.Server.Run(":3000")
}
