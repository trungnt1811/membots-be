package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const JWT_KEY = "JWT_PAYLOAD"
const USER_KEY = "AUTH_USER"

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

func CreateJWTMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwt.MapClaims{}

		parts := strings.Split(c.GetHeader("Authorization"), " ")
		tk := ""
		if len(parts) >= 2 {
			tk = parts[1]
		}
		if len(parts) == 1 {
			tk = parts[0]
		}

		if tk == "" {
			unauthorized(c, http.StatusUnauthorized, "authorization token required")
			return
		}

		jwt.ParseWithClaims(tk, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("mock verification key"), nil
		})

		// Check user
		switch v := claims["sub"].(type) {
		case nil:
			unauthorized(c, http.StatusBadRequest, "missing sub field")
			return
		case string:
			var user model.UserEntity
			result := db.First(&user, "id = ?", v)

			if result.Error != nil {
				if result.Error.Error() == "record not found" {
					unauthorized(c, http.StatusUnauthorized, "invalid user")
					return
				}
				unauthorized(c, http.StatusInternalServerError, "")
				return
			}

			c.Set(USER_KEY, user)
		default:
			unauthorized(c, http.StatusBadRequest, "sub must be string format")
			return
		}

		// Check
		switch v := claims["exp"].(type) {
		case nil:
			unauthorized(c, http.StatusBadRequest, "missing exp field")
			return
		case float64:
			if int64(v) < time.Now().Unix() {
				unauthorized(c, http.StatusUnauthorized, "token is expired")
				return
			}
		case json.Number:
			n, err := v.Int64()
			if err != nil {
				unauthorized(c, http.StatusBadRequest, "exp must be float64 format")
				return
			}
			if n < time.Now().Unix() {
				unauthorized(c, http.StatusUnauthorized, "token is expired")
				return
			}
		default:
			unauthorized(c, http.StatusBadRequest, "exp must be float64 format")
			return
		}

		c.Set(JWT_KEY, claims)
		c.Next()
	}
}

func GetJWTClaims(c *gin.Context) (jwt.MapClaims, error) {
	data, exist := c.Get(JWT_KEY)
	if !exist {
		return nil, errors.New("jwt payload not exist")
	}
	claims, ok := data.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("wrong jwt payload format")
	}
	return claims, nil
}

func GetAuthUser(c *gin.Context) (*model.UserEntity, error) {
	data, exist := c.Get(USER_KEY)
	if !exist {
		return nil, errors.New("auth user not exist")
	}
	user, ok := data.(model.UserEntity)
	if !ok {
		return nil, errors.New("wrong user entity payload format")
	}
	return &user, nil
}
