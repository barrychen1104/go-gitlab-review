package main

import (
	"go-gitlab-review/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	r.POST("/webhook", service.Webhook)

	r.Run(":5000")
}
