package handler

import (
	"core/config"
	"core/constant"
	"core/model"
	"core/request"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"time"
)

func QueryWeatherByLonLat(c *gin.Context) {
	// 解析请求坐标
	var position *request.QueryWeatherLonLatRequest
	if err := c.ShouldBind(&position); err != nil {
		ResponseFail(c, http.StatusBadRequest, constant.DataParseError, err.Error())
		return
	}

	// 将 Longitude 和 Latitude 转换为 float64 类型
	longitude, err := strconv.ParseFloat(position.Longitude, 64)
	if err != nil {
		ResponseFail(c, http.StatusBadRequest, constant.InvalidParameter, "经度参数格式错误")
		return
	}
	latitude, err := strconv.ParseFloat(position.Latitude, 64)
	if err != nil {
		ResponseFail(c, http.StatusBadRequest, constant.InvalidParameter, "纬度参数格式错误")
		return
	}

	// 对经纬度进行截断并格式化为两位小数
	position.Longitude = strconv.FormatFloat(math.Trunc(longitude*100)/100, 'f', 2, 64)
	position.Latitude = strconv.FormatFloat(math.Trunc(latitude*100)/100, 'f', 2, 64)
	// 向和风天气发起调用
	pos, err := QueryForPositionWithLonLat(position)
	if err != nil {
		ResponseFail(c, http.StatusInternalServerError, constant.CannotQueryPositionInfo, err.Error())
		return
	}
	weather, err := QueryForNowWeather(position)
	if err != nil {
		ResponseFail(c, http.StatusInternalServerError, constant.CannotQueryWeatherInfo, err.Error())
		return
	}
	airQuality, err := QueryAirQuality(position)
	if err != nil {
		ResponseFail(c, http.StatusInternalServerError, constant.CannotQueryAirQualityInfo, err.Error())
		return
	}
	city, ok := pos["name"].(string)
	if !ok {
		ResponseFail(c, http.StatusInternalServerError, constant.CannotQueryPositionInfo, "无法查询到城市名称")
		return
	}
	// 将 weather 转换为 JSON 字符串
	weatherJSON, err := json.Marshal(weather)
	if err != nil {
		ResponseFail(c, http.StatusInternalServerError, constant.CannotQueryWeatherInfo, "无法转换天气信息为 JSON 字符串")
		return
	}
	var queryRecord *model.QueryRecord
	queryRecord = &model.QueryRecord{
		City:        city,
		Longitude:   position.Longitude,
		Latitude:    position.Latitude,
		WeatherInfo: string(weatherJSON),
		CreateTime:  time.Now(),
		Type:        1,
		Deleted:     0,
	}
	if err := config.DataBase.Create(queryRecord).Error; err != nil {
		ResponseFail(c, http.StatusInternalServerError, constant.DataBaseSaveError, err.Error())
		return
	}
	data := map[string]interface{}{
		"pos":         pos,
		"weather":     weather,
		"air_quality": airQuality,
	}
	ResponseSuccessWithData(c, data)
}
