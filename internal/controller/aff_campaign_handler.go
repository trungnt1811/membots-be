package controller

import (
	"net/http"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/gin-gonic/gin"
)

type AffCampaignHandler struct {
	AffCampaignService interfaces.AffCampaignService
}

func NewAffCampaignHandler(
	affCampaignService interfaces.AffCampaignService,
) *AffCampaignHandler {
	return &AffCampaignHandler{
		AffCampaignService: affCampaignService,
	}
}

// GetAllAffCampaign Get list of all aff campaign
// @Summary Get list of all aff campaign
// @Description Get list of all aff campaign
// @Tags 	app
// @Accept	json
// @Produce json
// @Param page query string false "page to query, default is 1"
// @Param size query string false "size to query, default is 10"
// @Success 200 		{object}	dto.AffCampaignAppDtoResponse
// @Failure 401 		{object}	dto.GeneralError
// @Failure 400 		{object}	dto.GeneralError
// @Router 	/api/v1/app//aff-campaign [get]
func (handler *AffCampaignHandler) GetAllAffCampaign(ctx *gin.Context) {
	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	response, err := handler.AffCampaignService.GetAllAffCampaign(ctx, page, size)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "Get list of all aff campaign error: ", err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
