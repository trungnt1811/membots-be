package middleware

import (
	"fmt"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// NOTED: We use User Id as key => rate limit middleware must be place after authentication middleware
func keyFunc(c *gin.Context) string {
	user, err := dto.GetUserInfo(c)
	if err != nil {
		return c.FullPath() + c.ClientIP()
	}

	return c.FullPath() + fmt.Sprintf("%v", user.ID)
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

func NewRateLimit(redisClient *redis.Client, duration time.Duration, limit uint) gin.HandlerFunc {
	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: redisClient,
		Rate:        duration,
		Limit:       limit,
	})
	return ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})
}
