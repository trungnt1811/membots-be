package home_page

import (
	"net/http"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/gin-gonic/gin"
)

type HomePageHandler struct {
	HomePageUCase interfaces.HomePageUCase
}

func NewHomePageHandler(
	ucase interfaces.HomePageUCase,
) *HomePageHandler {
	return &HomePageHandler{
		HomePageUCase: ucase,
	}
}

// GetHomePage Get home page
// @Summary Get home page
// @Description Get home page
// @Tags app
// @Accept	json
// @Produce json
// @Success 200 		{object}	dto.HomePageDto
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

	response, err := handler.HomePageUCase.GetHomePage(ctx, uint64(user.ID))
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "get home page error", err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
