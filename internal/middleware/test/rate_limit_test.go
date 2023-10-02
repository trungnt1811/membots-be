package test

import (
	"fmt"
	"github.com/astraprotocol/affiliate-system/internal/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func dumpHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

func TestRateLimit(t *testing.T) {
	asserts := assert.New(t)

	limit := 10

	rdc := redis.NewClient(&redis.Options{})

	rateLimitMd := middleware.NewRateLimitMiddleware(func(c *gin.Context) (string, error) {
		return "unittest", nil
	}, "unittest", time.Duration(10*time.Second), rdc)

	// Setup server
	r := gin.Default()
	r.GET("/", rateLimitMd.WithGroup("unittest", limit), dumpHandler)

	// Run
	var lastResp *httptest.ResponseRecorder
	for i := 0; i < limit+1; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
		r.ServeHTTP(w, req)
		lastResp = w
	}
	// Expect last time to be error
	asserts.Equal(http.StatusTooManyRequests, lastResp.Code)
	fmt.Println(string(lastResp.Body.Bytes()))
}
