package middleware

import (
	"github.com/gin-gonic/gin"
)

type User struct {
	Name string
	Role string
}

func unauthorized(c *gin.Context, code int, message string) {
	c.Header("WWW-Authenticate", "JWT realm=affiliate-system")
	c.Abort()
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
