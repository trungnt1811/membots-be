package app_camp

import (
	"net/http"
	"strconv"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/gin-gonic/gin"
)

type AppCampHandler struct {
	AppCampService interfaces.AppCampService
}

func NewAppCampHandler(
	appCampService interfaces.AppCampService,
) *AppCampHandler {
	return &AppCampHandler{
		AppCampService: appCampService,
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
// @Router 	/api/v1/app/aff-campaign [get]
func (handler *AppCampHandler) GetAllAffCampaign(ctx *gin.Context) {
	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	response, err := handler.AppCampService.GetAllAffCampaign(ctx, page, size)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "Get list of all aff campaign error: ", err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// GetAffCampaignByAccesstradeId Get aff campaign by accesstrade id
// @Summary Get aff campaign by accesstrade id
// @Description Get aff campaign by accesstrade id
// @Tags 	app
// @Accept	json
// @Produce json
// @Param accesstradeId query string false "accesstradeId to query"
// @Success 200 		{object}	dto.AffCampaignAppDto
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/app/aff-campaign/{accresstradeId} [get]
func (handler *AppCampHandler) GetAffCampaignByAccesstradeId(ctx *gin.Context) {
	accesstradeId, err := strconv.Atoi(ctx.Param("accresstradeId"))
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "accresstradeId is invalid", err)
		return
	}

	response, err := handler.AppCampService.GetAffCampaignByAccesstradeId(ctx, uint64(accesstradeId))
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "Get list of all aff campaign error: ", err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
