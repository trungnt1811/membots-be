package aff_search

import (
	util "github.com/AstraProtocol/reward-libs/utils"
	"github.com/astraprotocol/affiliate-system/internal/interfaces"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AffSearchHandler struct {
	UCase interfaces.AffSearchUCase
}

func NewAffSearchHandler(service interfaces.AffSearchUCase) *AffSearchHandler {
	return &AffSearchHandler{
		UCase: service,
	}
}

// AffSearch search aff campaign
// @Summary search aff campaign
// @Description search aff campaign
// @Tags 	search
// @Accept	json
// @Produce json
// @Param q query string true "q to search"
// @Param page query string false "page to query, default is 1"
// @Success 200 		{object}	dto.AffSearchResponseDto
// @Failure 401 		{object}	util.GeneralError
// @Failure 400 		{object}	util.GeneralError
// @Router 	/api/v1/app/aff-search [get]
func (handler *AffSearchHandler) AffSearch(ctx *gin.Context) {
	q := ctx.DefaultQuery("q", "")
	page := ctx.GetInt("page")
	size := ctx.GetInt("size")
	response, err := handler.UCase.Search(ctx, q, page, size)
	if err != nil {
		util.RespondError(ctx, http.StatusInternalServerError, "search error", err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
