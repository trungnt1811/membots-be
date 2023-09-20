package statistic

import (
	"net/http"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/astraprotocol/affiliate-system/internal/util/log"
	"github.com/gin-gonic/gin"
)

type StatisticHandler struct {
	ucase interfaces.StatisticUcase
}

func NewStatisticHandler(ucase interfaces.StatisticUcase) *StatisticHandler {
	return &StatisticHandler{
		ucase: ucase,
	}
}

// GetSummary Get summary statistic data
// @Summary Get summary statistic data
// @Description Get summary statistic data
// @Tags 	console
// @Accept	json
// @Produce json
// @Security ApiKeyAuth
// @Query 	timeRange	query dto.TimeRange true "Time range"
// @Success 200 		{object}	dto.StatisticSummaryResponse
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/console/summary [get]
func (handler *StatisticHandler) GetSummary(ctx *gin.Context) {
	var d dto.TimeRange
	err := ctx.BindQuery(&d)
	if err != nil {
		log.LG.Errorf("parse time range error: %v", err)
	}

	resp, err := handler.ucase.GetSummaryByTimeRange(d)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "time range is required")
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
