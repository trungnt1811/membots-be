package user

import (
	"github.com/astraprotocol/affiliate-system/internal/middleware"
	"github.com/astraprotocol/affiliate-system/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetUserDetail Get a user information
// @Summary GetUserDetail
// @Description Get a user detail by uid
// @Tags user
// @Produce json
// @Param uid path string true "User uid to query"
// @Security ApiKeyAuth
// @Router /api/v1/users/{uid} [get]
func (handler *UserHandler) GetUserDetail(ctx *gin.Context) {
	jwtClaims, err := middleware.GetJWTClaims(ctx)
	if err != nil {
		util.RespondError(ctx, http.StatusBadRequest, "bad user claims", err)
		return
	}
	uid, ok := ctx.Params.Get("uid")
	if !ok {
		// Response error
		util.RespondError(ctx, http.StatusBadRequest, "bad user id", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"uid":    jwtClaims["sub"],
		"name":   "hello",
		"issuer": jwtClaims["iss"],
		"userId": uid,
	})
}
