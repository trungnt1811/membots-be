package campaign

import (
	"net/http"
	"strconv"

	util "github.com/AstraProtocol/reward-libs/utils"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	util2 "github.com/astraprotocol/affiliate-system/internal/util"

	"github.com/gin-gonic/gin"
)

type ConsoleBannerHandler struct {
	UCase interfaces.ConsoleBannerUCase
}

func NewConsoleBannerHandler(uCase interfaces.ConsoleBannerUCase) *ConsoleBannerHandler {
	return &ConsoleBannerHandler{
		UCase: uCase,
	}
}

// GetAllBanner Get list aff banner
// @Summary Get list aff banner
// @Description Get list aff banner
// @Tags console
// @Produce json
// @Param status query string false "by to query, default is all(active, inactive)"
// @Param order query string false "order to query, default is desc"
// @Param page query string false "page to query, default is 1"
// @Param size query string false "size to query, default is 10"
// @Success 200 		{object}	dto.AffBannerDtoResponse
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Security ApiKeyAuth
// @Router /api/v1/console/aff-banner [get]
func (handler *ConsoleBannerHandler) GetAllBanner(ctx *gin.Context) {
	queryStatus := ctx.DefaultQuery("status", "")
	order := ctx.DefaultQuery("order", "desc")

	queryStatusIn := util2.NormalizeStatusActiveInActive(queryStatus)
	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	if !util2.IsValidOrder(order) {
		util2.RespondError(ctx, http.StatusBadRequest, "query order invalid, order in {asc, desc}", nil)
		return
	}
	listAffBanner, err := handler.UCase.GetAllBanner(queryStatusIn, page, size)
	if err != nil {
		util2.RespondError(ctx, http.StatusInternalServerError, "get listBanner error", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, listAffBanner)
}

// UpdateBannerInfo update aff-banner info
// @Summary update aff-banner info
// @Description update aff-banner info
// @Tags 	console
// @Accept	json
// @Produce json
// @Param 	payload	body 			dto.AffBannerDto true "banner info to update, required"
// @Param id path int true "id to query"
// @Success 200 		{object}	dto.ResponseMessage
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Security ApiKeyAuth
// @Router 	/api/v1/console/aff-banner/{id} [PUT]
func (handler *ConsoleBannerHandler) UpdateBannerInfo(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "id is required", err)
		return
	}
	var payload dto.AffBannerDto
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "send payload is required", err)
		return
	}

	err = handler.UCase.UpdateBanner(uint(id), payload)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "failed to update aff banner", err)
		return
	}

	ctx.JSON(http.StatusOK, dto.ResponseMessage{
		Message: "success",
	})
}

// GetBannerById Get aff-banner by id
// @Summary Get aff-banner by id
// @Description Get aff-banner by id
// @Tags console
// @Produce json
// @Param id path int true "id to query"
// @Success 200 		{object}	dto.AffBannerDto
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Security ApiKeyAuth
// @Router /api/v1/console/aff-banner/{id} [get]
func (handler *ConsoleBannerHandler) GetBannerById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "id is required", err)
		return
	}

	affCampaign, err := handler.UCase.GetBannerById(uint(id))
	if err != nil {
		util2.RespondError(ctx, http.StatusInternalServerError, "get aff-banner error", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, affCampaign)
}
