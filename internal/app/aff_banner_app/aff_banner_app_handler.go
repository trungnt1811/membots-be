package campaign

import (
	"net/http"
	"strconv"

	util "github.com/AstraProtocol/reward-libs/utils"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	util2 "github.com/astraprotocol/affiliate-system/internal/util"

	"github.com/gin-gonic/gin"
)

type AppBannerHandler struct {
	UCase interfaces.AppBannerUCase
}

func NewAppBannerHandler(uCase interfaces.AppBannerUCase) *AppBannerHandler {
	return &AppBannerHandler{
		UCase: uCase,
	}
}

// GetAllBanner Get list aff banner
// @Summary Get list aff banner
// @Description Get list aff banner
// @Tags app-banner
// @Produce json
// @Param order query string false "order to query, default is desc"
// @Param page query string false "page to query, default is 1"
// @Param size query string false "size to query, default is 10"
// @Success 200 		{object}	dto.AffBannerDtoResponse
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Security ApiKeyAuth
// @Router /api/v1/app/aff-banner [get]
func (handler *AppBannerHandler) GetAllBanner(ctx *gin.Context) {
	order := ctx.DefaultQuery("order", "desc")
	queryStatus := "active"
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

// GetBannerById Get aff-banner by id
// @Summary Get aff-banner by id
// @Description Get aff-banner by id
// @Tags app-banner
// @Produce json
// @Param id path int true "id to query"
// @Success 200 		{object}	dto.AffBannerDto
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Security ApiKeyAuth
// @Router /api/v1/app/aff-banner/{id} [get]
func (handler *AppBannerHandler) GetBannerById(ctx *gin.Context) {
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
