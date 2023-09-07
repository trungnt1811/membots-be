package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CreateBasicAuthMiddleware(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		parts := strings.Split(c.GetHeader("Authorization"), " ")
		candidate := ""

		if len(parts) >= 2 {
			candidate = parts[1]
		}
		if len(parts) == 1 {
			candidate = parts[0]
		}

		if candidate == "" {
			unauthorized(c, http.StatusUnauthorized, "authorization token required")
			return
		}

		if candidate != token {
			unauthorized(c, http.StatusUnauthorized, "wrong token")
			return
		}

		c.Next()
	}
}
