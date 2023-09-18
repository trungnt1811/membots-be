package category

import (
	util "github.com/AstraProtocol/reward-libs/utils"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	util2 "github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AffCategoryHandler struct {
	UCase interfaces.AffCategoryUCase
}

func NewAffCategoryHandler(s interfaces.AffCategoryUCase) *AffCategoryHandler {
	return &AffCategoryHandler{
		UCase: s,
	}
}

// GetAllCategory Get all category have aff-campaign
// @Summary  Get all category have aff-campaign
// @Description Get all category have aff-campaign
// @Tags 	category
// @Accept	json
// @Produce json
// @Success 200 		{object}	dto.AffCategoryResponseDto
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/app/aff-categories [get]
func (handler *AffCategoryHandler) GetAllCategory(ctx *gin.Context) {
	page := ctx.GetInt("page")
	size := ctx.GetInt("size")
	response, err := handler.UCase.GetAllCategory(ctx, page, size)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "get all category error", err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

// GetAllAffCampaignInCategory Get all aff-campaign in category
// @Summary Get all aff-campaign in category
// @Description Get all aff-campaign in category
// @Tags 	category
// @Accept	json
// @Produce json
// @Param by query string false "by to query, default is ctime/price/top/sales"
// @Param order query string false "order to query, default is desc"
// @Success 200 		{object}	dtos.CouponDtoResponse
// @Failure 401 		{object}	dtos.GeneralError
// @Failure 400 		{object}	dtos.GeneralError
// @Param categoryId path int true "categoryId to query"
// @Param page query string false "page to query, default is 1"
// @Router 	/api/v1/categories/{categoryId} [get]
func (handler *AffCategoryHandler) GetAllAffCampaignInCategory(ctx *gin.Context) {
	categoryId, err := strconv.Atoi(ctx.Param("categoryId"))
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "categoryId id is required", err)
		return
	}
	page := ctx.GetInt("page")
	size := ctx.GetInt("size")
	queryBy := ctx.DefaultQuery("by", "ctime")
	order := ctx.DefaultQuery("order", "desc")
	if !util2.IsValidOrder(order) {
		util.RespondError(ctx, http.StatusBadRequest, "query order invalid, order in {asc, desc}", nil)
		return
	}
	response, err := handler.UCase.GetAllAffCampaignInCategory(ctx, uint32(categoryId), queryBy, order, page, size)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "Get all aff-campaign in category error ", err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
