package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/infra/caching"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req/v3"
)

type authHandler struct {
	HttpClient     *req.Client
	CreatorAuthUrl string
	AppAuthUrl     string
	RedisClient    caching.Repository
}

func (s *authHandler) CheckAdminHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token struct {
			Authorization *string `json:"Authorization" binding:"required"`
		}

		err := c.ShouldBindHeader(&token)
		if err != nil {
			util.RespondError(c, http.StatusUnauthorized, "no token provided", err)
			return
		}

		if len(*token.Authorization) < 200 {
			util.RespondError(c, http.StatusUnauthorized, "invalid token", err)
			return
		}

		info, err := s.creatorTokenInfo(*token.Authorization)
		if err != nil {
			log.Error().Msgf("check admin token error: %v", err)
			util.RespondError(c, http.StatusUnauthorized, err.Error())
			return
		}

		c.Set(dto.UserInfoKey, info)
		c.Next()
	}
}

func (s *authHandler) CheckUserHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token struct {
			Authorization *string `json:"Authorization" binding:"required"`
		}

		err := c.ShouldBindHeader(&token)
		if err != nil {
			util.RespondError(c, http.StatusUnauthorized, "no token provided", err)
			return
		}

		if len(*token.Authorization) < 200 {
			util.RespondError(c, http.StatusUnauthorized, "invalid token", err)
			return
		}

		info, err := s.appTokenInfo(*token.Authorization)
		if err != nil {
			log.Error().Msgf("check user token error: %v", err)
			util.RespondError(c, http.StatusUnauthorized, err.Error())
			return
		}

		c.Set(dto.UserInfoKey, info)
		c.Next()
	}
}

func (s *authHandler) creatorTokenInfo(jwtToken string) (dto.UserInfo, error) {
	var authInfo dto.UserCreatorResponse
	key := fmt.Sprint("creator_token_", jwtToken[200:len(jwtToken)-1])
	keyer := &caching.Keyer{Raw: key}
	err := s.RedisClient.RetrieveItem(keyer, &authInfo)
	if err != nil {
		userExpireInfo, err1 := s.getUserExpireFromJWTToken(jwtToken)
		if err1 != nil {
			return dto.UserInfo{}, err1
		}
		// cache miss
		resp, err1 := s.HttpClient.R().SetHeader("Authorization", jwtToken).
			SetSuccessResult(&authInfo).AddQueryParam("userId", fmt.Sprint(userExpireInfo.UserId)).
			Get(s.CreatorAuthUrl)
		if err1 != nil {
			return dto.UserInfo{}, err1
		}
		if !resp.IsSuccessState() {
			return dto.UserInfo{}, err1
		}
		expireAt := userExpireInfo.ExpiredAt - time.Now().Unix()
		if err = s.RedisClient.SaveItem(keyer, authInfo, time.Duration(expireAt)*time.Second); err != nil {
			return authInfo.Data[0].User, err
		}
	}
	return authInfo.Data[0].User, nil
}

func (s *authHandler) appTokenInfo(jwtToken string) (dto.UserInfo, error) {
	var authInfo dto.UserInfo
	key := fmt.Sprint("app_token_", jwtToken[200:len(jwtToken)-1])
	keyer := &caching.Keyer{Raw: key}
	err := s.RedisClient.RetrieveItem(keyer, &authInfo)
	if err != nil {
		userExpireInfo, err1 := s.getUserExpireFromJWTToken(jwtToken)
		if err1 != nil {
			return authInfo, err1
		}
		// cache miss
		resp, err1 := s.HttpClient.R().SetHeader("Authorization", jwtToken).
			SetSuccessResult(&authInfo).
			Get(s.AppAuthUrl)
		if err1 != nil {
			return dto.UserInfo{}, err1
		}
		if !resp.IsSuccessState() {
			return dto.UserInfo{}, err1
		}
		expireAt := userExpireInfo.ExpiredAt - time.Now().Unix()
		if err = s.RedisClient.SaveItem(keyer, authInfo, time.Duration(expireAt)*time.Second); err != nil {
			return authInfo, err
		}
	}
	return authInfo, nil
}

func (s *authHandler) getUserExpireFromJWTToken(jwtToken string) (dto.UserKeyExpiredDto, error) {
	claims := jwt.MapClaims{}
	var userKeyExpire dto.UserKeyExpiredDto

	if jwtToken == "" {
		return userKeyExpire, fmt.Errorf("jwtToken is empty")
	}

	jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return userKeyExpire, fmt.Errorf("mock verification key")
	})

	switch v := claims["sub"].(type) {
	case nil:
		return userKeyExpire, fmt.Errorf("missing sub field")
	case string:
		userId, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return userKeyExpire, fmt.Errorf("invalid user id")
		}
		userKeyExpire.UserId = uint32(userId)
	default:
		return userKeyExpire, fmt.Errorf("sub must be string format")
	}

	switch v := claims["exp"].(type) {
	case nil:
		return userKeyExpire, fmt.Errorf("missing exp field")
	case float64:
		if int64(v) < time.Now().Unix() {
			return userKeyExpire, fmt.Errorf("token is expired")
		}
		userKeyExpire.ExpiredAt = int64(v)
	case json.Number:
		n, err := v.Int64()
		if err != nil {
			return userKeyExpire, fmt.Errorf("exp must be float64 format")
		}
		if n < time.Now().Unix() {
			return userKeyExpire, fmt.Errorf("token is expired")
		}
		userKeyExpire.ExpiredAt = n
	default:
		return userKeyExpire, fmt.Errorf("exp must be float64 format")
	}
	return userKeyExpire, nil
}

func NewAuthUseCase(redisClient caching.Repository,
	creatorAuthUrl string,
	appAuthUrl string,
) *authHandler {
	client := req.C().
		SetUserAgent("affiliate-system"). // Chainable client settings.
		SetTimeout(3 * time.Second).
		SetCommonRetryCount(3).
		SetCommonErrorResult(&dto.ErrorMessage{}).
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if resp.Err != nil { // There is an underlying error, e.g. network error or unmarshal error.
				return nil
			}
			if errMsg, ok := resp.ErrorResult().(*dto.ErrorMessage); ok {
				resp.Err = errMsg // Convert api error into go error
				return nil
			}
			if !resp.IsSuccessState() {
				// Neither a success response nor a error response, record details to help troubleshooting
				resp.Err = fmt.Errorf("bad status: %s\nraw content:\n%s", resp.Status, resp.Dump())
			}
			return nil
		})
	return &authHandler{
		CreatorAuthUrl: creatorAuthUrl,
		AppAuthUrl:     appAuthUrl,
		HttpClient:     client,
		RedisClient:    redisClient,
	}
}
