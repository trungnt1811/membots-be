package aff_camp_app

import (
	"net/http"
	"strconv"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/gin-gonic/gin"
)

type AffCampAppHandler struct {
	AffCampAppService interfaces.AffCampAppService
}

func NewAffCampAppHandler(
	service interfaces.AffCampAppService,
) *AffCampAppHandler {
	return &AffCampAppHandler{
		AffCampAppService: service,
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
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/app/aff-campaign [get]
func (handler *AffCampAppHandler) GetAllAffCampaign(ctx *gin.Context) {
	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	response, err := handler.AffCampAppService.GetAllAffCampaign(ctx, page, size)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "Get list of all aff campaign error: ", err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// GetAffCampaignById Get aff campaign by id
// @Summary Get aff campaign by id
// @Description Get aff campaign by id
// @Tags 	app
// @Accept	json
// @Produce json
// @Param id path int true "id to query"
// @Success 200 		{object}	dto.AffCampaignAppDto
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/app/aff-campaign/{id} [get]
func (handler *AffCampAppHandler) GetAffCampaignById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "id is invalid", err)
		return
	}

	response, err := handler.AffCampAppService.GetAffCampaignById(ctx, uint64(id))
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "Get aff campaign by id error: ", err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
