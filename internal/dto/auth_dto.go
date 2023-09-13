package dto

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mitchellh/mapstructure"
)

const UserInfoKey = "user"

type UserDto struct {
	Id uint32 `json:"id"`
}

type UserKeyExpiredDto struct {
	UserId    uint32 `json:"user_id"`
	ExpiredAt int64  `json:"expired_at"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

func (msg *ErrorMessage) Error() string {
	return fmt.Sprintf("API Error: %s", msg.Message)
}

type JWTClaim struct {
	TokenType string   `json:"token_type"`
	Version   uint     `json:"version"`
	UserInfo  userInfo `json:"user"`
	Id        uint32   `json:"id"`
	Env       string   `json:"env"`
	jwt.RegisteredClaims
}

type userInfo struct {
	ID            uint32 `json:"id"`
	Email         string `json:"email,omitempty"`
	WalletAddress string `json:"wallet_address,omitempty"`
}

func GetUserInfo(ctx *gin.Context) (JWTClaim, error) {
	tmpUser, exists := ctx.Get(UserInfoKey)
	if !exists {
		return JWTClaim{}, fmt.Errorf("no token provided")
	}

	userJwt := JWTClaim{}
	err := mapstructure.Decode(tmpUser, &userJwt)
	if err != nil {
		return JWTClaim{}, fmt.Errorf("wrong form token")
	}
	return userJwt, nil
}
