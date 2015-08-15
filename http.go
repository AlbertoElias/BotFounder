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

	router.GET("/hello/:world", func(c *gin.Context) {
		name := c.Param("world")
		c.String(http.StatusOK, "Hello %s", name)
	})

	return &Server{router}
}

// Run starts the server
func (s *Server) Run() {
	s.Server.Run(":3000")
}
