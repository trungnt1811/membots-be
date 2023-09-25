package test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
)

const (
	WALLET_ADDR = "0xA4872a056C7cD4a6c60A880c8F3aCE886E05A8BF"
)

type UserTokenInfo struct {
	AccessToken  string         `json:"accessToken"`
	RefreshToken string         `json:"refreshToken"`
	User         map[string]any `json:"user"`
}

type AffRewardTestSuite struct {
	suite.Suite
	caller       *resty.Client
	JWTToken     string
	TestEndpoint string
	UserId       uint
}

func NewAffRewardTestSuite() *AffRewardTestSuite {
	client := resty.New()
	client.SetRetryCount(3)                                 // Retry 3 times
	client.SetTimeout(time.Duration(30 * int(time.Second))) // HTTP 30s timeout
	return &AffRewardTestSuite{
		caller:       client,
		TestEndpoint: "http://localhost:8080",
	}
}

func TestRunAffRewardTestSuite(t *testing.T) {
	suite.Run(t, NewAffRewardTestSuite())
}

// SetupSuite run setup before testing
func (s *AffRewardTestSuite) SetupSuite() {
	// First, setup account's JWT token
	req := s.caller.R()
	req.SetHeader("Content-Type", "application/json")
	req.SetBody(map[string]any{
		"address":      WALLET_ADDR,
		"expired_at":   1897834487,
		"created_at":   1687834893,
		"sign_message": "0x1689d82d5413952d80e3e43fc2efb1e10ce348716f9a9f3a0571057c28ee55592847cd62871e1aae36655d3c235aea7c468a72a2404f6f3e5b79be154008380601",
	})

	resp, err := req.Post("https://reward-app-backend.dev.astranet.services/api/v2/users")
	s.NoError(err)
	var userToken UserTokenInfo
	err = json.Unmarshal(resp.Body(), &userToken)
	s.NoError(err)

	s.JWTToken = userToken.AccessToken
	s.UserId = uint(userToken.User["id"].(float64))
}

func (s *AffRewardTestSuite) setDefaultHeader(req *resty.Request) {
	req.SetHeader("Authorization", s.JWTToken)
	req.SetHeader("Content-Type", "application/json")
}

func (s *AffRewardTestSuite) TestRunHappyCase() {
	// First, get generated aff link for this account
	req1 := s.caller.R()
	s.setDefaultHeader(req1)
	req1.SetBody(map[string]any{
		"campaign_id":  14, // Shopee
		"original_url": "",
		"shorten_link": false,
	})
	genLinkUrl := fmt.Sprint(s.TestEndpoint, "/api/v1/campaign/link")
	resp, err := req1.Post(genLinkUrl)
	s.NoError(err)
	s.True(resp.IsSuccess(), "http call status is not 200")
	var created dto.CreateLinkResponse
	err = json.Unmarshal(resp.Body(), &created)
	s.NoError(err)

	// Expect created link to have utm_content info
	createdUrl, err := url.Parse(created.AffLink)
	s.NoError(err)
	s.Contains(createdUrl.RawQuery, fmt.Sprintf("utm_content=%d-", s.UserId), "utm_content is missing")

	// Expect

	// After that, send a mock order post back for this test

	// After order initialized, update order approved with post back

	// Then, try to request ASA cashback for the order

	// AFter cashback run, make sure balance is updated
}
