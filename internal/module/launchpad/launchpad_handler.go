package launchpad

import (
	"fmt"
	"net/http"

	util "github.com/AstraProtocol/reward-libs/utils"
	"github.com/gin-gonic/gin"

	"github.com/flexstack.ai/membots-be/internal/interfaces"
)

type LaunchpadHandler struct {
	UCase interfaces.LaunchpadUCase
}

func NewLaunchpadHandler(ucase interfaces.LaunchpadUCase) *LaunchpadHandler {
	return &LaunchpadHandler{
		UCase: ucase,
	}
}

// GetHistoryByAddress Get launchpad history by address
// @Summary Get launchpad history by address
// @Description Get launchpad history by address
// @Tags 	launchpad
// @Accept	json
// @Produce json
// @Param memeAddress query string false "memeAddress to query, default is "
// @Success 200 		{object}	dto.LaunchpadInfoResp
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/launchpad [get]
func (handler *LaunchpadHandler) GetHistoryByAddress(ctx *gin.Context) {
	memeAddress := ctx.DefaultQuery("memeAddress", "")

	if memeAddress == "" {
		util.RespondError(ctx, http.StatusInternalServerError, "Get launchpad by memeAddress error: ", fmt.Errorf("memeAddress is empty"))
		return
	}

	response, err := handler.UCase.GetHistory(ctx, memeAddress)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "Get launchpad by memeAddress error: ", err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
