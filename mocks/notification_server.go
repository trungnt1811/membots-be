package mocks

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Mocking Notification Services API for debugging purpose
func RunMockServer() {
	r := gin.Default()
	r.Any("*", func(c *gin.Context) {
		// Response with success
		c.JSON(http.StatusOK, gin.H{
			"body": c.Request.Body,
		})
	})

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
