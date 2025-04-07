package handler

import (
	"compress/gzip"
	"core/config"
	"core/constant"
	"core/request"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func QueryForPositionWithLonLat(position *request.QueryWeatherLonLatRequest) (map[string]interface{}, error) {
	strLon, strLat := position.Longitude, position.Latitude
	queryCurPosUrl := config.UrbanSearchAPI + strLon + "," + strLat
	req, err := http.NewRequest("GET", queryCurPosUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(config.APIKeyHeader, config.ApiKey)
	// 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	// 解析响应
	if resp.StatusCode != http.StatusOK {
		// 修复：打印状态码和请求 URL
		fmt.Printf("Request failed with status code: %d, URL: %s\n", resp.StatusCode, resp.Request.URL.String())
		return nil, err
	}
	var reader io.ReadCloser = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
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
		return nil, err
	}
	// 提取第一个位置信息
	locationInterface, ok := response["location"]
	if !ok {
		return nil, constant.NewError(-1, constant.CannotQueryPositionInfo, "无法解析位置信息")
	}
	locations, ok := locationInterface.([]interface{})
	if !ok {
		return nil, constant.NewError(-1, constant.CannotQueryPositionInfo, "无法解析位置信息")
	}
	if len(locations) == 0 {
		return nil, constant.NewError(-1, constant.CannotQueryPositionInfo, "无法解析位置信息")
	}
	firstLocationInterface, ok := locations[0].(map[string]interface{})
	if !ok {
		return nil, constant.NewError(-1, constant.CannotQueryPositionInfo, "无法解析位置信息")
	}
	return firstLocationInterface, nil
}

func QueryForNowWeather(position *request.QueryWeatherLonLatRequest) (map[string]interface{}, error) {
	strLon, strLat := position.Longitude, position.Latitude
	queryNowWeatherUrl := config.RealTimeWeatherQueryAPI + strLon + "," + strLat
	req, err := http.NewRequest("GET", queryNowWeatherUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(config.APIKeyHeader, config.ApiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, constant.NewError(-1, constant.CannotQueryWeatherInfo, "无法查询到天气信息")
	}
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	// 提取now
	nowInterface, ok := response["now"]
	if !ok {
		return nil, constant.NewError(-1, constant.CannotQueryWeatherInfo, "无法查询到天气信息")
	}
	now, ok := nowInterface.(map[string]interface{})
	if !ok {
		return nil, constant.NewError(-1, constant.CannotQueryWeatherInfo, "无法查询到天气信息")
	}
	return now, nil
}

func QueryAirQuality(position *request.QueryWeatherLonLatRequest) (map[string]interface{}, error) {
	strLon, strLat := position.Longitude, position.Latitude
	queryAirQualityUrl := config.AirQualityQueryAPI + strLat + "/" + strLon
	req, err := http.NewRequest("GET", queryAirQualityUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(config.APIKeyHeader, config.ApiKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status code: %d, URL: %s\n", resp.StatusCode, resp.Request.URL.String())
		return nil, constant.NewError(-1, constant.CannotQueryAirQualityInfo, "无法查询到天气信息")
	}
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	// 提取 indexes
	indexesInterface, ok := response["indexes"]
	if !ok {
		return nil, constant.NewError(-1, constant.CannotQueryAirQualityInfo, "无法查询到天气信息")
	}
	indexes, ok := indexesInterface.([]interface{})
	if !ok {
		return nil, constant.NewError(-1, constant.CannotQueryAirQualityInfo, "无法查询到天气信息")
	}
	if len(indexes) == 0 {
		return nil, constant.NewError(-1, constant.CannotQueryAirQualityInfo, "无法查询到天气信息")
	}
	firstIndexInterface, ok := indexes[0].(map[string]interface{})
	if !ok {
		return nil, constant.NewError(-1, constant.CannotQueryAirQualityInfo, "无法查询到天气信息")
	}
	return firstIndexInterface, nil
}
