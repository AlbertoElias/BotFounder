package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// State takes care of storing global dependencies of the project
var State = struct {
	*DB
}{}

func main() {

	db, err := SetupDb()
	panicOnErr(err)
	State.DB = db

	router := gin.Default()

	router.GET("/hello/:world", func(c *gin.Context) {
		name := c.Param("world")
		c.String(http.StatusOK, "Hello %s", name)
	})

	router.Run(":3000")
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
