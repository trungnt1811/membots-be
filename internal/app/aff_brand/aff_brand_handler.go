package aff_brand

import (
	"net/http"

	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/gin-gonic/gin"
)

type AffBrandHandler struct {
	AffBrandUCase interfaces.AffBrandUCase
}

func NewAffBrandHandler(
	ucase interfaces.AffBrandUCase,
) *AffBrandHandler {
	return &AffBrandHandler{
		AffBrandUCase: ucase,
	}
}

// GetTopFavouriteAffBrand Get top favorited aff brands
// @Summary Get top favorited aff brands
// @Description Get top favorited aff brands
// @Tags 	app
// @Accept	json
// @Produce json
// @Param page query string false "page to query, default is 1"
// @Param size query string false "size to query, default is 10"
// @Success 200 		{object}	[]dto.AffCampaignLessDto
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Security ApiKeyAuth
// @Router 	/api/v1/app/brand/top-favorited [get]
func (handler *AffBrandHandler) GetTopFavouriteAffBrand(ctx *gin.Context) {
	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	response, err := handler.AffBrandUCase.GetTopFavouriteAffBrand(ctx, page, size)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "Get top favorited aff brands error: ", err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
