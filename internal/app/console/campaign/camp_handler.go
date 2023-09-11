package campaign

import (
	util "github.com/AstraProtocol/reward-libs/utils"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	util2 "github.com/astraprotocol/affiliate-system/internal/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ConsoleCampHandler struct {
	UCase interfaces.ConsoleCampUCase
}

func NewConsoleCampHandler(uCase interfaces.ConsoleCampUCase) *ConsoleCampHandler {
	return &ConsoleCampHandler{
		UCase: uCase,
	}
}

// GetAllCampaign Get list aff campaign
// @Summary Get list aff campaign
// @Description Get list aff campaign
// @Tags console
// @Produce json
// @Param stella_status query string false "by to query, default is all"
// @Param order query string false "order to query, default is desc"
// @Param page query string false "page to query, default is 1"
// @Param size query string false "size to query, default is 10"
// @Success 200 		{object}	dto.AffCampaignDtoResponse
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router /api/v1/console/aff-campaign [get]
func (handler *ConsoleCampHandler) GetAllCampaign(ctx *gin.Context) {
	queryStatus := ctx.DefaultQuery("stella_status", "")
	order := ctx.DefaultQuery("order", "desc")

	queryStatusIn := util2.NormalizeStatus(queryStatus)
	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	if !util2.IsValidOrder(order) {
		util2.RespondError(ctx, http.StatusBadRequest, "query order invalid, order in {asc, desc}", nil)
		return
	}
	listAffCampaign, err := handler.UCase.GetAllCampaign(queryStatusIn, page, size)
	if err != nil {
		util2.RespondError(ctx, http.StatusInternalServerError, "get listAffCampaign error", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, listAffCampaign)
}

// UpdateCampaignInfo update campaign info
// @Summary update campaign info
// @Description update campaign info
// @Tags 	console
// @Accept	json
// @Produce json
// @Param 	payload	body 			dto.AffCampaignAppDto true "Campaign info to update, required"
// @Param id path int true "id to query"
// @Success 200 		{object}	dto.AffCampaignAppDto
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Security ApiKeyAuth
// @Router 	/api/v1/console/aff-campaign/{id} [PUT]
func (handler *ConsoleCampHandler) UpdateCampaignInfo(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "id is required", err)
		return
	}
	var payload dto.AffCampaignAppDto
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "send payload is required", err)
		return
	}

	err = handler.UCase.UpdateCampaign(uint(id), payload)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "failed to update user info", err)
		return
	}

	ctx.JSON(http.StatusOK, "success")
}

// GetCampaignById Get campaign by id
// @Summary Get campaign by id
// @Description Get campaign by id
// @Tags console
// @Produce json
// @Param id path int true "id to query"
// @Success 200 		{object}	dto.AffCampaignDto
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router /api/v1/console/aff-campaign/{id} [get]
func (handler *ConsoleCampHandler) GetCampaignById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "id is required", err)
		return
	}

	affCampaign, err := handler.UCase.GetCampaignById(uint(id))
	if err != nil {
		util2.RespondError(ctx, http.StatusInternalServerError, "get aff-campaign error", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, affCampaign)
}
