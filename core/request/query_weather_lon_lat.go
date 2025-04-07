package request

type QueryWeatherLonLatRequest struct {
	Longitude string `form:"longitude" json:"longitude" binding:"required"`
	Latitude  string `form:"latitude" json:"latitude" binding:"required"`
}
