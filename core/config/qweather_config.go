package config

var (
	SCHEMA       = "https"
	ApiHost      = "kv4bj974qq.re.qweatherapi.com"
	APIKeyHeader = "X-QW-Api-Key"
	ApiKey       = "9569aefb17be426a9e4a0fec738e8aac"

	UrbanSearchAPIPath          = "/geo/v2/city/lookup"
	RealTimeWeatherQueryAPIPath = "/v7/weather/now"
	AirQualityQueryAPIPath      = "/airquality/v1/current/"

	UrbanSearchAPI          = SCHEMA + "://" + ApiHost + UrbanSearchAPIPath + "?location="
	RealTimeWeatherQueryAPI = SCHEMA + "://" + ApiHost + RealTimeWeatherQueryAPIPath + "?location="
	AirQualityQueryAPI      = SCHEMA + "://" + ApiHost + AirQualityQueryAPIPath
)
