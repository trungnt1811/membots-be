package user_view_aff_camp

import (
	"net/http"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/gin-gonic/gin"
)

type UserViewAffCampHandler struct {
	UserViewAffCampUCase interfaces.UserViewAffCampUCase
}

func NewUserViewAffCampHandler(
	ucase interfaces.UserViewAffCampUCase,
) *UserViewAffCampHandler {
	return &UserViewAffCampHandler{
		UserViewAffCampUCase: ucase,
	}
}

// GetListRecentlyVisitedSection Get list recently visited section by user
// @Summary Get list recently visited by user
// @Description Get list recently visited by user
// @Tags 	app
// @Accept	json
// @Produce json
// @Param page query string false "page to query, default is 1"
// @Param size query string false "size to query, default is 10"
// @Success 200 		{object}	dto.AffCampaignComBrandDtoResponse
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Security ApiKeyAuth
// @Router 	/api/v1/app/recently-visited-section [get]
func (handler *UserViewAffCampHandler) GetListRecentlyVisitedSection(ctx *gin.Context) {
	// First, take user from JWT
	user, err := dto.GetUserInfo(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "logged in user required", err)
		return
	}

	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	response, err := handler.UserViewAffCampUCase.GetListUserViewAffCampByUserId(ctx, uint64(user.ID), page, size)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "Get list of recently visited section error: ", err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
