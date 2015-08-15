package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	db, err := SetupDb()
	panicOnErr(err)
	fmt.Println("hey")

	fmt.Println(db)

	router := gin.Default()

	router.GET("/hello/:world", func(c *gin.Context) {
		name := c.Param("world")
		c.String(http.StatusOK, "Hello %s", name)
	})

	router.Run(":3000")
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
