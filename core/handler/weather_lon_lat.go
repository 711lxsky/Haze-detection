package handler

import (
	"core/config"
	"core/constant"
	"core/model"
	"core/request"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
	"time"
)

func QueryWeatherByLonLat(c *gin.Context) {
	// 解析请求坐标
	var position request.QueryWeatherLonLatRequest
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
	pos, err := QueryForPositionWithLonLat(&position)
	if err != nil {
		ResponseFail(c, http.StatusInternalServerError, constant.CannotQueryPositionInfo, err.Error())
		return
	}
	weather, err := QueryForNowWeather(&position)
	if err != nil {
		ResponseFail(c, http.StatusInternalServerError, constant.CannotQueryWeatherInfo, err.Error())
		return
	}
	airQuality, err := QueryAirQuality(&position)
	if err != nil {
		ResponseFail(c, http.StatusInternalServerError, constant.CannotQueryAirQualityInfo, err.Error())
		return
	}
	city := pos.Name
	//if !ok {
	//	ResponseFail(c, http.StatusInternalServerError, constant.CannotQueryPositionInfo, "无法查询到城市名称")
	//	return
	//}
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
	fmt.Printf("pos: %+v\n", pos)
	fmt.Printf("weather: %+v\n", weather)
	fmt.Printf("airQuality: %+v\n\n", airQuality)
	var data = Data{
		Pos:        pos,
		Weather:    weather,
		AirQuality: airQuality,
	}
	ResponseSuccessWithData(c, data)
}

// 定义返回的数据结构
type ResponseData struct {
	Data    Data   `json:"data"`
	Message string `json:"message"`
}

type Data struct {
	AirQuality AirQuality `json:"air_quality"`
	Pos        Position   `json:"pos"`
	Weather    Weather    `json:"weather"`
}

type AirQuality struct {
	Aqi              int       `json:"aqi"`
	AqiDisplay       string    `json:"aqiDisplay"`
	Category         string    `json:"category"`
	Code             string    `json:"code"`
	Color            Color     `json:"color"`
	Health           Health    `json:"health"`
	PrimaryPollutant Pollutant `json:"primaryPollutant"`
}

type Color struct {
	Alpha float64 `json:"alpha"`
	Blue  int     `json:"blue"`
	Green int     `json:"green"`
	Red   int     `json:"red"`
}

type Health struct {
	Advice Advice `json:"advice"`
	Level  string `json:"level"`
	Name   string `json:"name"`
}

type Advice struct {
	GeneralPopulation   string `json:"generalPopulation"`
	SensitivePopulation string `json:"sensitivePopulation"`
	Effect              string `json:"effect"`
}

type Pollutant struct {
	Code     string `json:"code"`
	FullName string `json:"fullName"`
	Name     string `json:"name"`
}

type Position struct {
	Adm1      string `json:"adm1"`
	Aweather  string `json:"aweather.com/weather/chang'an-101110102.html"`
	Id        string `json:"id"`
	IsDst     string `json:"isDst"`
	Lat       string `json:"lat"`
	Lon       string `json:"lon"`
	Name      string `json:"name"`
	Rank      string `json:"rank"`
	Type      string `json:"type"`
	Tz        string `json:"tz"`
	UtcOffset string `json:"utcOffset"`
}

type Weather struct {
	Cloud           string `json:"cloud"`
	Dew             string `json:"dew"`
	FlsLike         string `json:"flsLike"`
	Humidity        string `json:"humidity"`
	Icon            string `json:"icon"`
	ObservationTime string `json:"obsTime"`
	Precip          string `json:"precip"`
	Pressure        string `json:"pressure"`
	Temp            string `json:"temp"`
	Text            string `json:"text"`
	Vis             string `json:"vis"`
	Wind360         string `json:"wind360"`
	WindDir         string `json:"windDir"`
	WindScale       string `json:"windScale"`
	WindSpeed       string `json:"windSpeed"`
}
