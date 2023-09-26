package test

import (
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"os"
	"testing"
	"time"

	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade"
	"github.com/astraprotocol/affiliate-system/internal/app/accesstrade/types"
	logger "github.com/astraprotocol/affiliate-system/internal/util/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

const (
	TEST_APIKEY = "wCITgIrAXX53MQ3uo8Z28-evlUY1Lcqn"
)

type RepositoryTestSuite struct {
	suite.Suite
	repository interfaces.ATRepository
}

func NewRepositoryTestSuite() *RepositoryTestSuite {
	logger.LG = logger.NewZerologLogger(os.Stdout, zerolog.InfoLevel)
	repository := accesstrade.NewAccessTradeRepository(TEST_APIKEY, 3, 30)
	return &RepositoryTestSuite{
		repository: repository,
	}
}

func TestRunRepositorySuite(t *testing.T) {
	suite.Run(t, NewRepositoryTestSuite())
}

// SetupSuite run setup before testing
func (s *RepositoryTestSuite) SetupSuite() {
}

func (s *RepositoryTestSuite) TestQueryCampaigns() {
	page := 1
	resp, err := s.repository.QueryCampaigns(false, page, 10)
	s.NoError(err)
	s.Equal(page, resp.Page)
}

func (s *RepositoryTestSuite) TestQueryOrders() {
	page := 1
	q := types.ATOrderQuery{
		Since: time.Date(2022, 8, 1, 0, 0, 0, 0, time.UTC),
		Until: time.Date(2022, 8, 31, 0, 0, 0, 0, time.UTC),
	}
	resp, err := s.repository.QueryOrders(q, page, 10)
	s.NoError(err)
	// Expect to be at least 1 order item
	s.Equal(page, resp.Total)
}

func (s *RepositoryTestSuite) TestCreateLink() {
	page := 1
	resp, err := s.repository.QueryCampaigns(true, page, 10)
	s.NoError(err)
	s.GreaterOrEqual(len(resp.Data), 1)
	camp := resp.Data[0]
	url := camp.Url
	s.NotEmpty(url)

	linkResp, err := s.repository.CreateTrackingLinks(camp.Id, true, []string{url}, map[string]string{
		"utm_campaign": "stella",
		"utm_source":   "testing",
	})
	s.NoError(err)
	s.GreaterOrEqual(len(linkResp.Data.SuccessLink), 1)
}

func (s *RepositoryTestSuite) TestCreateLinkFailed() {
	campaignId := "4348614231480407268" // Tikivn
	url := "http://shopee.vn"           // Shopee

	linkResp, err := s.repository.CreateTrackingLinks(campaignId, true, []string{url}, map[string]string{
		"utm_campaign": "stella",
		"utm_source":   "testing",
	})
	s.ErrorContains(err, "The link is not part of the campaign")
	s.Nil(linkResp)
}
