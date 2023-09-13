package dto

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mitchellh/mapstructure"
)

const UserInfoKey = "user"
const UserCreatorInfoKey = "user_creator"
const UserAppInfoKey = "user_app"

type UserDto struct {
	Id uint32 `json:"id"`
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
	UserInfo  UserInfo `json:"user"`
	Id        uint32   `json:"id"`
	Env       string   `json:"env"`
	jwt.RegisteredClaims
}

type UserInfo struct {
	ID            uint32 `json:"id"`
	Email         string `json:"email,omitempty"`
	WalletAddress string `json:"wallet_address,omitempty"`
}

type UserCreatorResponse struct {
	Data []UserCreatorData `json:"data"`
}

type UserCreatorData struct {
	UserId uint64   `json:"userId"`
	RoleId uint64   `json:"roleId,omitempty"`
	User   UserInfo `json:"user"`
}

func GetUserId(ctx *gin.Context) (uint32, error) {
	tmpUser, exists := ctx.Get(UserInfoKey)
	if !exists {
		return 0, fmt.Errorf("no token provided")
	}
	userDto := UserDto{}
	err := mapstructure.Decode(tmpUser, &userDto)
	if err != nil {
		return 0, fmt.Errorf("wrong form token")
	}
	return userDto.Id, nil
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

func GetUserCreatorInfo(ctx *gin.Context) (UserInfo, error) {
	tmpUser, exists := ctx.Get(UserCreatorInfoKey)
	if !exists {
		return UserInfo{}, fmt.Errorf("no token provided")
	}
	userInfo := UserInfo{}
	err := mapstructure.Decode(tmpUser, &userInfo)
	if err != nil {
		return UserInfo{}, fmt.Errorf("cannot decode user creator info")
	}
	return userInfo, nil
}
