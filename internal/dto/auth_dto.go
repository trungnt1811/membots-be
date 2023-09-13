package dto

import (
	"fmt"

	"github.com/gin-gonic/gin"
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
	userInfo := UserInfo{}
	err := mapstructure.Decode(tmpUser, &userInfo)
	if err != nil {
		return 0, fmt.Errorf("cannot decode user info")
	}
	return userInfo.ID, nil
}

func GetUserInfo(ctx *gin.Context) (UserInfo, error) {
	tmpUser, exists := ctx.Get(UserInfoKey)
	if !exists {
		return UserInfo{}, fmt.Errorf("no token provided")
	}
	userInfo := UserInfo{}
	err := mapstructure.Decode(tmpUser, &userInfo)
	if err != nil {
		return UserInfo{}, fmt.Errorf("cannot decode user info")
	}
	return userInfo, nil
}
