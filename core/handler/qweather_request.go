package handler

import (
	"compress/gzip"
	"core/config"
	"core/constant"
	"core/request"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"io"
	"net/http"
)

func QueryForPositionWithLonLat(position *request.QueryWeatherLonLatRequest) (Position, error) {
	strLon, strLat := position.Longitude, position.Latitude
	queryCurPosUrl := config.UrbanSearchAPI + strLon + "," + strLat
	req, err := http.NewRequest("GET", queryCurPosUrl, nil)
	if err != nil {
		return Position{}, err
	}
	req.Header.Set(config.APIKeyHeader, config.ApiKey)
	// 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Position{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	// 解析响应
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status code: %d, URL: %s\n", resp.StatusCode, resp.Request.URL.String())
		return Position{}, constant.NewError(-1, constant.CannotQueryPositionInfo, "无法查询到位置信息")
	}
	var reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return Position{}, err
		}
		defer func(reader io.ReadCloser) {
			err := reader.Close()
			if err != nil {
				return
			}
		}(reader)
	}
	var response map[string]interface{}
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		return Position{}, err
	}
	// 提取第一个位置信息
	locationInterface, ok := response["location"]
	if !ok {
		return Position{}, constant.NewError(-1, constant.CannotQueryPositionInfo, "无法解析位置信息")
	}
	locations, ok := locationInterface.([]interface{})
	if !ok {
		return Position{}, constant.NewError(-1, constant.CannotQueryPositionInfo, "无法解析位置信息")
	}
	if len(locations) == 0 {
		return Position{}, constant.NewError(-1, constant.CannotQueryPositionInfo, "无法解析位置信息")
	}
	firstLocationInterface, ok := locations[0].(map[string]interface{})
	if !ok {
		return Position{}, constant.NewError(-1, constant.CannotQueryPositionInfo, "无法解析位置信息")
	}
	var pos Position
	if err := mapstructure.Decode(firstLocationInterface, &pos); err != nil {
		return Position{}, err
	}
	return pos, nil
}

func QueryForNowWeather(position *request.QueryWeatherLonLatRequest) (Weather, error) {
	strLon, strLat := position.Longitude, position.Latitude
	queryNowWeatherUrl := config.RealTimeWeatherQueryAPI + strLon + "," + strLat
	req, err := http.NewRequest("GET", queryNowWeatherUrl, nil)
	if err != nil {
		return Weather{}, err
	}
	req.Header.Set(config.APIKeyHeader, config.ApiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Weather{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return Weather{}, constant.NewError(-1, constant.CannotQueryWeatherInfo, "无法查询到天气信息")
	}
	var reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return Weather{}, err
		}
		defer func(reader io.ReadCloser) {
			err := reader.Close()
			if err != nil {
				return
			}
		}(reader)
	}
	var response map[string]interface{}
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		return Weather{}, err
	}
	// 提取now
	nowInterface, ok := response["now"]
	if !ok {
		return Weather{}, constant.NewError(-1, constant.CannotQueryWeatherInfo, "无法查询到天气信息")
	}
	now, ok := nowInterface.(map[string]interface{})
	if !ok {
		return Weather{}, constant.NewError(-1, constant.CannotQueryWeatherInfo, "无法查询到天气信息")
	}
	var weather Weather
	if err := mapstructure.Decode(now, &weather); err != nil {
		return Weather{}, err
	}
	return weather, nil
}

func QueryAirQuality(position *request.QueryWeatherLonLatRequest) (AirQuality, error) {
	strLon, strLat := position.Longitude, position.Latitude
	queryAirQualityUrl := config.AirQualityQueryAPI + strLat + "/" + strLon
	req, err := http.NewRequest("GET", queryAirQualityUrl, nil)
	if err != nil {
		return AirQuality{}, err
	}
	req.Header.Set(config.APIKeyHeader, config.ApiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return AirQuality{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status code: %d, URL: %s\n", resp.StatusCode, resp.Request.URL.String())
		return AirQuality{}, constant.NewError(-1, constant.CannotQueryAirQualityInfo, "无法查询到天气信息")
	}
	var reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return AirQuality{}, err
		}
		defer func(reader io.ReadCloser) {
			err := reader.Close()
			if err != nil {
				return
			}
		}(reader)
	}
	var response map[string]interface{}
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		return AirQuality{}, err
	}
	// 提取 indexes
	indexesInterface, ok := response["indexes"]
	if !ok {
		return AirQuality{}, constant.NewError(-1, constant.CannotQueryAirQualityInfo, "无法查询到天气信息")
	}
	indexes, ok := indexesInterface.([]interface{})
	if !ok {
		return AirQuality{}, constant.NewError(-1, constant.CannotQueryAirQualityInfo, "无法查询到天气信息")
	}
	if len(indexes) == 0 {
		return AirQuality{}, constant.NewError(-1, constant.CannotQueryAirQualityInfo, "无法查询到天气信息")
	}
	firstIndexInterface, ok := indexes[0].(map[string]interface{})
	if !ok {
		return AirQuality{}, constant.NewError(-1, constant.CannotQueryAirQualityInfo, "无法查询到天气信息")
	}
	var airQuality AirQuality
	if err := mapstructure.Decode(firstIndexInterface, &airQuality); err != nil {
		return AirQuality{}, err
	}
	return airQuality, nil
}
