package test

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
)

const (
	API_ENDPOINT = "http://localhost:8080/api/v1"
	WALLET_ADDR  = "0xA4872a056C7cD4a6c60A880c8F3aCE886E05A8BF"
	RPC_ENDPOINT = "https://rpc.astranaut.dev"
)

type UserTokenInfo struct {
	AccessToken  string         `json:"accessToken"`
	RefreshToken string         `json:"refreshToken"`
	User         map[string]any `json:"user"`
}

type AffRewardTestSuite struct {
	suite.Suite
	caller      *resty.Client
	JWTToken    string
	UserId      uint
	UserWallet  common.Address
	UserBalance *big.Int
	RPCClient   *ethclient.Client
}

func NewAffRewardTestSuite() *AffRewardTestSuite {
	client := resty.New()
	client.SetRetryCount(3)                                 // Retry 3 times
	client.SetTimeout(time.Duration(30 * int(time.Second))) // HTTP 30s timeout
	return &AffRewardTestSuite{
		caller: client,
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
	client, err := ethclient.Dial(RPC_ENDPOINT)
	s.NoError(err)
	s.RPCClient = client

	s.UserWallet = common.HexToAddress(WALLET_ADDR[2:])

	currentBlock, err := client.BlockNumber(context.Background())
	s.NoError(err)
	balance, err := client.BalanceAt(context.Background(), s.UserWallet, big.NewInt(int64(currentBlock)))
	s.NoError(err)
	s.True(balance.Cmp(big.NewInt(0)) != 0)
	s.UserBalance = balance
}

func (s *AffRewardTestSuite) setDefaultHeader(req *resty.Request) {
	req.SetHeader("Authorization", s.JWTToken)
	req.SetHeader("Content-Type", "application/json")
}

func getPart(s string, sep string, idx int) string {
	parts := strings.Split(s, sep)
	if idx >= len(parts) {
		return ""
	}
	return parts[idx]
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
	genLinkUrl := fmt.Sprint(API_ENDPOINT, "/campaign/link")
	resp, err := req1.Post(genLinkUrl)
	s.NoError(err)
	s.True(resp.IsSuccess(), "http call status is not 200")
	var created dto.CreateLinkResponse
	err = json.Unmarshal(resp.Body(), &created)
	s.NoError(err)

	// Expect created link to have utm_content info
	createdUrl, err := url.Parse(created.AffLink)
	s.NoError(err)
	s.Contains(createdUrl.RawQuery, fmt.Sprintf("utm_content=%d-", s.UserId), "utm_content is missing user id")

	// Get utm_content from string
	idx := strings.LastIndex(created.AffLink, "utm_content=")
	utmContent := getPart(created.AffLink[idx:], "=", 1)
	trackedId := getPart(utmContent, "-", 1)
	s.NotEmpty(utmContent)
	s.NotEmpty(trackedId)

	// After that, send a mock order post back for update order approved
	req2 := s.caller.R()
	s.setDefaultHeader(req2)
	req1.SetBody(map[string]any{
		"transaction_id":       "199759877",
		"campaign_id":          "4751584435713464237",
		"order_id":             "230918N260WKSG",
		"product_id":           "15688627684@shopee@20562",
		"quantity":             1,
		"product_category":     "Cameras_&_Flycam",
		"product_price":        20562.0,
		"reward":               206.0,
		"sales_time":           "2023-09-18 08:11:45.000000",
		"click_time":           "2023-09-18 08:09:18.000000",
		"browser":              "Mobile Safari",
		"conversion_platform":  "mobile_app",
		"ip":                   "52.77.0.178",
		"referrer":             "",
		"utm_source":           "stella",
		"utm_campaign":         "",
		"utm_content":          utmContent,
		"utm_medium":           "",
		"status":               0,
		"publisher_login_name": "astrarewards",
		"is_confirmed":         0,
		"customer_type":        ""})
	postBackUrl := fmt.Sprint(API_ENDPOINT, "/order/post-back")
	resp, err = req1.Post(postBackUrl)
	s.NoError(err)
	s.True(resp.IsSuccess(), "http call status is not 200")
	s.Contains(string(resp.Body()), "230918N260WKSG") // Contains order id in resp

	// Sleep and wait for order reward to be created
	time.Sleep(time.Duration(time.Second * 5))

	// Then, try to request ASA cashback for the order
	req3 := s.caller.R()
	s.setDefaultHeader(req3)
	req3.SetBody(map[string]any{})
	withdrawUrl := fmt.Sprint(API_ENDPOINT, "/app/rewards/withdraw")
	resp, err = req3.Post(withdrawUrl)
	s.NoError(err)
	s.True(resp.IsSuccess())
	var withdrawResp dto.RewardWithdrawDto
	err = json.Unmarshal(resp.Body(), &withdrawResp)
	s.NoError(err)
	s.Condition(func() bool {
		return withdrawResp.Amount > 0
	}, "withdraw amount is 0")

	// Sleep and wait for withdraw reward to be executed
	time.Sleep(time.Duration(time.Second * 5))

	// After withdraw processed, make sure balance is updated
	currentBlk, err := s.RPCClient.BlockNumber(context.Background())
	afterBalance, err := s.RPCClient.BalanceAt(context.Background(), s.UserWallet, big.NewInt(int64(currentBlk)))
	s.Condition(func() bool {
		return afterBalance.Cmp(s.UserBalance) > 0
	}, "after withdraw, balance not change")
}
