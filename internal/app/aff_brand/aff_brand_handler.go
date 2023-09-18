package aff_brand

import (
	"net/http"
	"strconv"

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
// @Param topFavorite query string false "topFavorite to query, default is 7 (1 <= topFavorite <= 10)"
// @Success 200 		{object}	[]dto.AffCampaignComBrandDto
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Security ApiKeyAuth
// @Router 	/api/v1/app/brand/top-favorited [get]
func (handler *AffBrandHandler) GetTopFavouriteAffBrand(ctx *gin.Context) {
	topFavoriteParam := ctx.DefaultQuery("topFavorite", "7")
	topFavorite, err := strconv.Atoi(topFavoriteParam)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "topFavorite must be integer")
		return
	}
	if topFavorite < 1 || topFavorite > 10 {
		util.RespondError(ctx, http.StatusBadRequest, "topFavorite must be 1 to 10")
		return
	}

	response, err := handler.AffBrandUCase.GetTopFavouriteAffBrand(ctx, topFavorite)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "Get top favrorited aff brands error: ", err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
