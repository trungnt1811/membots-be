package statistic

import (
	"net/http"
	"strconv"

	"github.com/astraprotocol/affiliate-system/internal/dto"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"github.com/gin-gonic/gin"
)

type StatisticHandler struct {
	ucase interfaces.StatisticUCase
}

func NewStatisticHandler(ucase interfaces.StatisticUCase) *StatisticHandler {
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
// @Param 	timeRange	query dto.TimeRange false "Time range to summary"
// @Success 200 		{object}	dto.StatisticSummaryResponse
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/console/summary [get]
func (handler *StatisticHandler) GetSummary(ctx *gin.Context) {
	var d dto.TimeRange
	err := ctx.BindQuery(&d)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "parse time range failed", err)
		return
	}

	resp, err := handler.ucase.GetSummaryByTimeRange(d)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "get summary failed")
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetCampaignSummary Get campaign statistic data
// @Summary Get campaign statistic data
// @Description Get campaign statistic data
// @Tags 	console
// @Accept	json
// @Produce json
// @Security ApiKeyAuth
// @Param 	campaignId	path uint 					true "Time range"
// @Param 	timeRange	  query dto.TimeRange false "Time range"
// @Success 200 		{object}	dto.CampaignSummaryResponse
// @Failure 424 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/console/summary/:campaignId [get]
func (handler *StatisticHandler) GetCampaignSummary(ctx *gin.Context) {
	var d dto.TimeRange
	err := ctx.BindQuery(&d)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "parse time range failed", err)
		return
	}
	campaignIdParam, _ := ctx.Params.Get("campaignId")
	campaignId, err := strconv.ParseUint(campaignIdParam, 10, 64)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "parse campaignId failed", err)
		return
	}
	if campaignId == 0 {
		util.RespondError(ctx, http.StatusBadRequest, "campaignId is required")
		return
	}

	resp, err := handler.ucase.GetCampaignSummaryByTimeRange(uint(campaignId), d)
	if err != nil {
		util.RespondError(ctx, http.StatusFailedDependency, "get summary failed")
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
