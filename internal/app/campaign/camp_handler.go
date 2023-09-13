package campaign

import (
	"net/http"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/gin-gonic/gin"
)

type CampaignHandler struct {
	usecase interfaces.CampaignUCase
}

func NewCampaignHandler(usecase interfaces.CampaignUCase) *CampaignHandler {
	return &CampaignHandler{
		usecase: usecase,
	}
}

// PostGenerateAffLink Request create a new campaign link or pick the active old one
// @Summary PostGenerateAffLink
// @Description Request create a new campaign link or pick the active old one
// @Tags 		redeem
// @Accept	json
// @Produce json
// @Param 	payload	body 			dto.CreateLinkPayload true "Request create link payload, required"
// @Success 200 		{object}	dto.CreateLinkResponse "when success, return the created link for this request campaign"
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Security ApiKeyAuth
// @Router 	/api/v1/campaign/link [post]
func (handler *CampaignHandler) PostGenerateAffLink(c *gin.Context) {
	// First, take user from JWT
	user, err := dto.GetUserInfo(c)
	if err != nil {
		util.RespondError(c, http.StatusBadRequest, "logged in user required", err)
		return
	}
	// Then verify payload data
	var payload dto.CreateLinkPayload
	err = c.BindJSON(&payload)
	if err != nil {
		util.RespondError(c, http.StatusBadRequest, "payload required", err)
		return
	}

	link, err := handler.usecase.GenerateAffLink(uint64(user.UserInfo.ID), &payload)
	if err != nil {
		util.RespondError(c, http.StatusFailedDependency, "create link fail", err)
		return
	}

	// Response transaction status
	c.JSON(http.StatusOK, link)
}
