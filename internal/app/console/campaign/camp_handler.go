package campaign

import (
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	util2 "github.com/astraprotocol/affiliate-system/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConsoleCampHandler struct {
	UCase interfaces.ConsoleCampUCase
}

func NewConsoleCampHandler(uCase interfaces.ConsoleCampUCase) *ConsoleCampHandler {
	return &ConsoleCampHandler{
		UCase: uCase,
	}
}

// GetAllCampaign Get list aff campaign
// @Summary Get list aff campaign
// @Description Get list aff campaign
// @Tags console
// @Produce json
// @Param stella_status query string false "by to query, default is all"
// @Param order query string false "order to query, default is desc"
// @Param page query string false "page to query, default is 1"
// @Param size query string false "size to query, default is 10"
// @Success 200 		{object}	dto.AffCampaignDtoResponse
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router /api/v1/console/aff-campaign [get]
func (handler *ConsoleCampHandler) GetAllCampaign(ctx *gin.Context) {
	queryStatus := ctx.DefaultQuery("stella_status", "")
	order := ctx.DefaultQuery("order", "desc")

	queryStatusIn := util2.NormalizeStatus(queryStatus)
	page := ctx.GetInt("page")
	size := ctx.GetInt("size")

	if !util2.IsValidOrder(order) {
		util2.RespondError(ctx, http.StatusBadRequest, "query order invalid, order in {asc, desc}", nil)
		return
	}
	listAffCampaign, err := handler.UCase.GetAllCampaign(queryStatusIn, page, size)
	if err != nil {
		util2.RespondError(ctx, http.StatusInternalServerError, "get listAffCampaign error", err)
		return
	}

	// Response transaction status
	ctx.JSON(http.StatusOK, listAffCampaign)
}
