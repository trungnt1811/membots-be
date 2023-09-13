package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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
		// cache miss
		userId, err1 := s.getUserIdFromJWTToken(jwtToken)
		if err1 != nil {
			return dto.UserInfo{}, err1
		}
		resp, err1 := s.HttpClient.R().SetHeader("Authorization", jwtToken).
			SetSuccessResult(&authInfo).AddQueryParam("userId", fmt.Sprint(userId)).
			Get(s.CreatorAuthUrl)
		if err1 != nil {
			return dto.UserInfo{}, resp.Err
		}
		if !resp.IsSuccessState() {
			return dto.UserInfo{}, resp.Err
		}
		if err = s.RedisClient.SaveItem(keyer, authInfo, time.Minute); err != nil {
			return dto.UserInfo{}, err
		}
	}
	return authInfo.Data[0].User, nil
}

func (s *authHandler) appTokenInfo(jwtToken string) (dto.UserDto, error) {
	var authInfo dto.UserDto
	key := fmt.Sprint("app_token_", jwtToken[200:len(jwtToken)-1])
	keyer := &caching.Keyer{Raw: key}
	err := s.RedisClient.RetrieveItem(keyer, &authInfo)
	if err != nil {
		// cache miss
		resp, err1 := s.HttpClient.R().SetHeader("Authorization", jwtToken).
			SetSuccessResult(&authInfo).
			Get(s.AppAuthUrl)
		if err1 != nil {
			return authInfo, resp.Err
		}
		if !resp.IsSuccessState() {
			return authInfo, resp.Err
		}
		if err = s.RedisClient.SaveItem(keyer, authInfo, time.Minute); err != nil {
			return authInfo, err
		}
	}
	return authInfo, nil
}

func (s *authHandler) getUserIdFromJWTToken(jwtToken string) (uint64, error) {
	claims := jwt.MapClaims{}

	if jwtToken == "" {
		return 0, fmt.Errorf("jwtToken is empty")
	}

	jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return 0, fmt.Errorf("mock verification key")
	})

	switch v := claims["sub"].(type) {
	case nil:
		return 0, fmt.Errorf("missing sub field")
	case string:
		userId, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid user id")
		}
		return userId, nil
	default:
		return 0, fmt.Errorf("sub must be string format")
	}
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
