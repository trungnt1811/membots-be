package test

import (
	"fmt"
	"os"
	"testing"

	atMocks "github.com/astraprotocol/affiliate-system/internal/app/accesstrade/mocks"
	"github.com/astraprotocol/affiliate-system/internal/app/campaign"
	"github.com/astraprotocol/affiliate-system/internal/app/campaign/mocks"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/model"
	logger "github.com/astraprotocol/affiliate-system/internal/util/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type CampaignUcaseTestSuite struct {
	suite.Suite
	ucase *campaign.CampaignUsecase
}

func NewCampaignUcaseTestSuite() *CampaignUcaseTestSuite {
	logger.LG = logger.NewZerologLogger(os.Stdout, zerolog.InfoLevel)
	atRepo := atMocks.NewAccessTradeRepoMock()
	campRepo := mocks.NewCampaignRepoMock()
	ucase := campaign.NewCampaignUsecase(campRepo, atRepo)
	return &CampaignUcaseTestSuite{
		ucase: ucase,
	}
}

func TestCampaignUcaseSuiteRun(t *testing.T) {
	suite.Run(t, NewCampaignUcaseTestSuite())
}

func (s *CampaignUcaseTestSuite) TestGenerateAffLink() {
	payload := &dto.CreateLinkPayload{
		CampaignId:  14,
		OriginalUrl: "",
		ShortenLink: false,
	}
	user := &model.UserEntity{
		ID: 2449,
	}
	resp, err := s.ucase.GenerateAffLink(user, payload)
	s.NoError(err)

	s.Equal(payload.CampaignId, resp.CampaignId)
	s.Contains(resp.AffLink, fmt.Sprintf("utm_content=%d", user.ID))
}
