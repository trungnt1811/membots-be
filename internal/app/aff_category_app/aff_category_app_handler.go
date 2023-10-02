package category

import (
	"fmt"
	util "github.com/AstraProtocol/reward-libs/utils"
	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
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

// GetListAffBrandByUser Get list aff brand by category
// @Summary Get list aff brand by category
// @Description Get list aff brand by category
// @Tags 	category
// @Accept	json
// @Produce json
// @Param categoryId path int true "categoryId to query"
// @Param page query string false "page to query, default is 1"
// @Param size query string false "size to query, default is 10"
// @Param filter query string false "filter to query, default is top-favorited (top-favorited/most-commission)"
// @Success 200 		{object}	dto.AffCampaignAppDtoResponse
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Security ApiKeyAuth
// @Router 	/api/v1/app/aff-categories/{categoryId} [get]
func (handler *AffCategoryHandler) GetListAffBrandByUser(ctx *gin.Context) {
	categoryId, err := strconv.Atoi(ctx.Param("categoryId"))
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "categoryId id is required", err)
		return
	}
	// First, take user from JWT
	user, err := dto.GetUserInfo(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
		return
	}

	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	filter := ctx.DefaultQuery("filter", "top-favorited")
	switch filter {
	case "top-favorited":
		response, err := handler.UCase.GetTopFavouriteAffBrand(ctx, uint(categoryId), uint64(user.ID), page, size)
		if err != nil {
			util.RespondError(ctx, http.StatusInternalServerError, "Get top favorited aff brands error: ", err)
			return
		}
		ctx.JSON(http.StatusOK, response)
		return
	case "most-commission":
		response, err := handler.UCase.GetMostCommissionAffCampaign(ctx, uint(categoryId), uint64(user.ID), page, size)
		if err != nil {
			util.RespondError(ctx, http.StatusInternalServerError, "Get most commission aff brands error: ", err)
			return
		}
		ctx.JSON(http.StatusOK, response)
		return
	default:
		util.RespondError(ctx, http.StatusBadRequest, fmt.Sprintf("filter: %s is not supported", filter))
		return
	}

}
