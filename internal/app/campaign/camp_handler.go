package campaign

import (
	"net/http"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/gin-gonic/gin"
)

type CampaignHandler struct {
	uCase interfaces.CampaignUCase
}

func NewCampaignHandler(uCase interfaces.CampaignUCase) *CampaignHandler {
	return &CampaignHandler{
		uCase: uCase,
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
func (handler *CampaignHandler) PostGenerateAffLink(ctx *gin.Context) {
	// First, take user from JWT
	user, err := dto.GetUserInfo(ctx)
	// The commented code is checking if there is an error while getting the user information from the JWT
	// token. If there is an error, it will respond with a HTTP status code of 400 (Bad Request) and an
	// error message indicating that a logged in user is required.
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
		return
	}
	// Then verify payload data
	var payload dto.CreateLinkPayload
	err = ctx.BindJSON(&payload)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "payload required", err)
		return
	}

	link, err := handler.uCase.GenerateAffLink(uint64(user.ID), &payload)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "create link fail", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, link)
}
