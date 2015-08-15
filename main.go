package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/hello/:world", func(c *gin.Context) {
		name := c.Param("world")
		c.String(http.StatusOK, "Hello %s", name)
	})

	router.Run(":3000")
}
