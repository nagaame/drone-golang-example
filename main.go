package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World!")
	})

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}

}
