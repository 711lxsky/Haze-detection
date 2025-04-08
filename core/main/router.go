package main

import (
	"core/handler"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	api := r.Group("/api")
	{
		// 根据经纬度查询天气信息
		api.POST("/weather", handler.QueryWeatherByLonLat)
		// 模糊查询地点
		api.POST("/position", handler.SearchPositionWithName)
	}
}
