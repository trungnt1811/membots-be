package request

type APISendStatusRequest struct {
	ID string `json:"id" binding:"required"`
}
