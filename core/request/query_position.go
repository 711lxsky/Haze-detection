package request

type QueryPositionRequest struct {
	Position string `form:"position" json:"position" binding:"required"`
}
