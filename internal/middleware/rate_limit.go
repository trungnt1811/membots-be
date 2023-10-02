package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func hitRateLimit(c *gin.Context, message string) {
	c.Abort()
	c.JSON(http.StatusTooManyRequests, gin.H{
		"code":    http.StatusTooManyRequests,
		"message": message,
	})
}

type IdFunc func(c *gin.Context) (string, error)

type RateLimitMiddleware struct {
	IdFunc    IdFunc
	GroupName string
	Duration  time.Duration
	Rdc       *redis.Client
}

func NewRateLimitMiddleware(idFunc IdFunc,
	groupName string,
	duration time.Duration,
	rdc *redis.Client) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		IdFunc:    idFunc,
		GroupName: groupName,
		Duration:  duration,
		Rdc:       rdc,
	}
}

func (md *RateLimitMiddleware) WithGroup(groupName string, limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := md.IdFunc(c)
		if err != nil {
			unauthorized(c, http.StatusBadRequest, "require id")
		}

		key := fmt.Sprintf("%s:%s", groupName, id)

		ctx := context.Background()
		cmd := md.Rdc.Get(ctx, key)
		err = cmd.Err()
		if err != nil {
			if err.Error() == "redis: nil" {
				// Set the item
				setCmd := md.Rdc.Set(ctx, key, 1, md.Duration)
				err = setCmd.Err()
				if err != nil {
					c.Abort()
					c.JSON(http.StatusBadGateway, gin.H{
						"message": "cannot set rate limit",
					})
				}
			}
			c.Abort()
			c.JSON(http.StatusBadGateway, gin.H{
				"message": "cannot get rate limit",
			})
			return
		}
		// If not error, try to get int
		num, err := cmd.Int()
		if err != nil {
			c.Abort()
			c.JSON(http.StatusBadGateway, gin.H{
				"message": "cannot parse rate limit",
			})
			return
		}
		intCmd := md.Rdc.Incr(ctx, key)
		err = intCmd.Err()
		if err != nil {
			c.Abort()
			c.JSON(http.StatusBadGateway, gin.H{
				"message": "cannot incr rate limit",
			})
			return
		}

		if num < limit {
			c.Next()
		} else {
			// Hit limit, throw rate limit error
			hitRateLimit(c, "too many requests")
		}
	}
}
