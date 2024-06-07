package stats

import (
	"fmt"
	"net/http"

	util "github.com/AstraProtocol/reward-libs/utils"
	"github.com/gin-gonic/gin"

	"github.com/flexstack.ai/membots-be/internal/interfaces"
)

type StatsHandler struct {
	UCase interfaces.StatsUCase
}

func NewStatsHandler(ucase interfaces.StatsUCase) *StatsHandler {
	return &StatsHandler{
		UCase: ucase,
	}
}

// GetStatsByMemeAddress Get stats by memeAddress
// @Summary Get stats by memeAddress
// @Description Get stats by memeAddress
// @Tags 	stats
// @Accept	json
// @Produce json
// @Param memeAddress query string false "memeAddress to query"
// @Success 200 		{object}	interface{}
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/stats [get]
func (handler *StatsHandler) GetStatsByMemeAddress(ctx *gin.Context) {
	memeAddress := ctx.DefaultQuery("memeAddress", "")
	url := fmt.Sprintf("https://www.trugly.meme/api/v1/stats?memeAddress=%s", memeAddress)
	response, err := handler.UCase.GetStatsByMemeAddress(ctx, url)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "Get starts error: ", err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
