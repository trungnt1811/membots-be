package auth

import (
	"fmt"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/caching"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/dto"
	"github.com/astraprotocol/affiliate-system/internal/app/affiliate/util"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req/v3"
	"net/http"
	"time"
)

type AuthHandler struct {
	HttpClient     *req.Client
	CreatorAuthUrl string
	AppAuthUrl     string
	RedisClient    caching.Repository
}

func (s *AuthHandler) CheckAdminHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token struct {
			Authorization *string `json:"Authorization" binding:"required"`
		}

		err := c.ShouldBindHeader(&token)
		if err != nil {
			c.Next()
			return
		}

		if len(*token.Authorization) < 300 {
			c.Next()
			return
		}

		info, err := s.creatorTokenInfo(*token.Authorization)
		if err != nil {
			fmt.Println(err)
			util.RespondError(c, http.StatusUnauthorized, err.Error())
			return
		}

		c.Set(dto.UserInfoKey, info)
		c.Next()
	}
}

func (s *AuthHandler) CheckUserHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		var token struct {
			Authorization *string `json:"Authorization" binding:"required"`
		}

		err := c.ShouldBindHeader(&token)
		if err != nil {
			c.Next()
			return
		}

		if len(*token.Authorization) < 300 {
			c.Next()
			return
		}

		info, err := s.appTokenInfo(*token.Authorization)
		if err != nil {
			fmt.Println(err)
			util.RespondError(c, http.StatusUnauthorized, err.Error())
			return
		}

		c.Set(dto.UserInfoKey, info)
		c.Next()
	}
}

func (s *AuthHandler) creatorTokenInfo(jwtToken string) (dto.UserDto, error) {
	var authInfo dto.UserDto
	key := fmt.Sprint("token_", jwtToken[200:len(jwtToken)-1])
	keyer := &caching.Keyer{Raw: key}
	err := s.RedisClient.RetrieveItem(keyer, &authInfo)
	if err != nil {
		// cache miss
		resp, err1 := s.HttpClient.R().SetHeader("Authorization", jwtToken).
			SetSuccessResult(&authInfo).
			Get(s.CreatorAuthUrl)
		if err1 != nil {
			return authInfo, err1
		}
		if !resp.IsSuccessState() {
			return authInfo, err1
		}
		if err = s.RedisClient.SaveItem(keyer, authInfo, time.Minute); err != nil {
			return authInfo, err
		}
	}
	return authInfo, nil
}

func (s *AuthHandler) appTokenInfo(jwtToken string) (dto.UserDto, error) {
	var authInfo dto.UserDto
	key := fmt.Sprint("token_", jwtToken[200:len(jwtToken)-1])
	keyer := &caching.Keyer{Raw: key}
	err := s.RedisClient.RetrieveItem(keyer, &authInfo)
	if err != nil {
		// cache miss
		resp, err1 := s.HttpClient.R().SetHeader("Authorization", jwtToken).
			SetSuccessResult(&authInfo).
			Get(s.CreatorAuthUrl)
		if err1 != nil {
			return authInfo, err1
		}
		if !resp.IsSuccessState() {
			return authInfo, err1
		}
		if err = s.RedisClient.SaveItem(keyer, authInfo, time.Minute); err != nil {
			return authInfo, err
		}
	}
	return authInfo, nil
}

func NewAuthUseCase(redisClient caching.Repository,
	creatorAuthUrl string,
	appAuthUrl string,
) *AuthHandler {
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
	return &AuthHandler{
		CreatorAuthUrl: creatorAuthUrl,
		AppAuthUrl:     appAuthUrl,
		HttpClient:     client,
		RedisClient:    redisClient,
	}
}
