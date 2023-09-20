package aff_brand

import (
	"fmt"
	"net/http"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/gin-gonic/gin"
)

type AffBrandHandler struct {
	UserViewAffCampUCase interfaces.UserViewAffCampUCase
	AffBrandUCase        interfaces.AffBrandUCase
}

func NewAffBrandHandler(
	userViewAffCampUCase interfaces.UserViewAffCampUCase,
	affBrandUCase interfaces.AffBrandUCase,
) *AffBrandHandler {
	return &AffBrandHandler{
		UserViewAffCampUCase: userViewAffCampUCase,
		AffBrandUCase:        affBrandUCase,
	}
}

// GetListAffBrandByUser Get list aff brand by user
// @Summary Get list aff brand by user
// @Description Get list aff brand by user
// @Tags 	app
// @Accept	json
// @Produce json
// @Param page query string false "page to query, default is 1"
// @Param size query string false "size to query, default is 10"
// @Param filter query string false "filter to query, default is recently-visited (recently-visited/top-favorited/favorite)"
// @Success 200 		{object}	dto.AffCampaignAppDtoResponse
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Security ApiKeyAuth
// @Router 	/api/v1/app/brand [get]
func (handler *AffBrandHandler) GetListAffBrandByUser(ctx *gin.Context) {
	// First, take user from JWT
	user, err := dto.GetUserInfo(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
		return
	}

	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	filter := ctx.DefaultQuery("filter", "recently-visited")
	switch filter {
	case "recently-visited":
		response, err := handler.UserViewAffCampUCase.GetListUserViewAffCampByUserId(ctx, uint64(user.ID), page, size)
		if err != nil {
			util.RespondError(ctx, http.StatusInternalServerError, "Get list of recently visited section error: ", err)
			return
		}
		ctx.JSON(http.StatusOK, response)
		return
	case "favorite":
		response, err := handler.AffBrandUCase.GetListFavAffBrandByUserId(ctx, uint64(user.ID), page, size)
		if err != nil {
			util.RespondError(ctx, http.StatusInternalServerError, "Get list fav brand by user error: ", err)
			return
		}
		ctx.JSON(http.StatusOK, response)
		return
	case "top-favorited":
		response, err := handler.AffBrandUCase.GetTopFavouriteAffBrand(ctx, uint64(user.ID), page, size)
		if err != nil {
			util.RespondError(ctx, http.StatusInternalServerError, "Get top favorited aff brands error: ", err)
			return
		}
		ctx.JSON(http.StatusOK, response)
		return
	default:
		util.RespondError(ctx, http.StatusBadRequest, fmt.Sprintf("filter: %s is not supported", filter))
		return
	}

}
