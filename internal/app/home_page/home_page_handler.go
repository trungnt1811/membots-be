package home_page

import (
	"net/http"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/gin-gonic/gin"
)

type HomePageHandler struct {
	AffBrandUCase interfaces.AffBrandUCase
}

func NewHomePageHandler(
	ucase interfaces.AffBrandUCase,
) *HomePageHandler {
	return &HomePageHandler{
		AffBrandUCase: ucase,
	}
}

// GetHomePage Get home page
// @Summary Get home page
// @Description Get home page
// @Tags app
// @Accept	json
// @Produce json
// @Success 200 		{object}	dto.AffCampaignAppDtoResponse
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Security ApiKeyAuth
// @Router /api/v1/app/home-page [get]
func (handler *HomePageHandler) GetHomePage(ctx *gin.Context) {
	// First, take user from JWT
	user, err := dto.GetUserInfo(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
		return
	}

	response, err := handler.AffBrandUCase.GetTopFavouriteAffBrand(ctx, uint64(user.ID), 1, 15)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "get home page error", err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
