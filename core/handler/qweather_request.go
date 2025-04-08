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
	"time"
)

// APIHandler 封装所有API请求逻辑
type APIHandler struct {
	client *http.Client
}

// NewAPIHandler 创建带有优化配置的API处理器
func NewAPIHandler() *APIHandler {
	return &APIHandler{
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    60 * time.Second,
				DisableCompression: false,
			},
		},
	}
}

// -------------------- 公共方法 --------------------

// doRequest 统一处理HTTP请求
func (h *APIHandler) doRequest(url string) (map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("请求创建失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set(config.APIKeyHeader, config.ApiKey)

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求执行失败: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, constant.NewError(
			-1,
			constant.RequestErr,
			fmt.Sprintf("状态码: %d, 响应: %s", resp.StatusCode, string(bodyBytes)),
		)
	}

	return parseJSONResponse(resp)
}

// parseJSONResponse 统一处理响应解析
func parseJSONResponse(resp *http.Response) (map[string]interface{}, error) {
	reader, err := getResponseReader(resp)
	if err != nil {
		return nil, err
	}
	defer closeReader(reader)

	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w, 原始数据: %s", err, string(bodyBytes))
	}

	return response, nil
}

// getResponseReader 处理内容编码
func getResponseReader(resp *http.Response) (io.ReadCloser, error) {
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("gzip解压失败: %w", err)
		}
		return gzReader, nil
	default:
		return resp.Body, nil
	}
}

// closeReader 安全关闭可关闭的读取器
func closeReader(reader io.ReadCloser) {
	if closer, ok := reader.(io.Closer); ok {
		_ = closer.Close()
	}
}

// -------------------- 业务函数 --------------------

// QueryForPositionWithLonLat 查询位置信息
func (h *APIHandler) QueryForPositionWithLonLat(position *request.QueryWeatherLonLatRequest) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s,%s", config.UrbanSearchAPI, position.Longitude, position.Latitude)
	response, err := h.doRequest(url)
	if err != nil {
		return nil, fmt.Errorf("位置查询失败: %w", err)
	}

	return extractLocation(response)
}

func extractLocation(response map[string]interface{}) (map[string]interface{}, error) {
	location, ok := response["location"]
	if !ok {
		return nil, constant.NewError(
			-1,
			constant.ResponseErr,
			"响应缺少location字段",
		)
	}

	locations, ok := location.([]interface{})
	if !ok || len(locations) == 0 {
		return nil, constant.NewError(
			-1,
			constant.ResponseErr,
			"location格式无效",
		)
	}

	firstLocation, ok := locations[0].(map[string]interface{})
	if !ok {
		return nil, constant.NewError(
			-1,
			constant.ResponseErr,
			"首位置信息格式无效",
		)
	}

	return firstLocation, nil
}

// QueryForNowWeather 查询实时天气
func (h *APIHandler) QueryForNowWeather(position *request.QueryWeatherLonLatRequest) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s,%s", config.RealTimeWeatherQueryAPI, position.Longitude, position.Latitude)
	response, err := h.doRequest(url)
	if err != nil {
		return nil, fmt.Errorf("实时天气查询失败: %w", err)
	}

	return extractNowWeather(response)
}

func extractNowWeather(response map[string]interface{}) (map[string]interface{}, error) {
	now, ok := response["now"]
	if !ok {
		return nil, constant.NewError(
			-1,
			constant.ResponseErr,
			"响应缺少now字段",
		)
	}

	nowMap, ok := now.(map[string]interface{})
	if !ok {
		return nil, constant.NewError(
			-1,
			constant.ResponseErr,
			"now字段格式无效",
		)
	}

	return nowMap, nil
}

// QueryAirQuality 查询空气质量
func (h *APIHandler) QueryAirQuality(position *request.QueryWeatherLonLatRequest) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s/%s", config.AirQualityQueryAPI, position.Latitude, position.Longitude)
	response, err := h.doRequest(url)
	if err != nil {
		return nil, fmt.Errorf("空气质量查询失败: %w", err)
	}

	return extractAirQuality(response)
}

func extractAirQuality(response map[string]interface{}) (map[string]interface{}, error) {
	indexes, ok := response["indexes"]
	if !ok {
		return nil, constant.NewError(
			-1,
			constant.ResponseErr,
			"响应缺少indexes字段",
		)
	}

	indexList, ok := indexes.([]interface{})
	if !ok || len(indexList) == 0 {
		return nil, constant.NewError(
			-1,
			constant.ResponseErr,
			"indexes格式无效",
		)
	}

	firstIndex, ok := indexList[0].(map[string]interface{})
	if !ok {
		return nil, constant.NewError(
			-1,
			constant.ResponseErr,
			"首项空气质量数据格式无效",
		)
	}

	return firstIndex, nil
}

// QueryNextWeather 查询未来天气
func (h *APIHandler) QueryNextWeather(position *request.QueryWeatherLonLatRequest) ([]interface{}, error) {
	url := fmt.Sprintf("%s%s,%s", config.QueryNextWeatherAPI, position.Longitude, position.Latitude)
	response, err := h.doRequest(url)
	if err != nil {
		return nil, fmt.Errorf("未来天气查询失败: %w", err)
	}

	return extractDailyWeather(response)
}

func extractDailyWeather(response map[string]interface{}) ([]interface{}, error) {
	daily, ok := response["daily"]
	if !ok {
		return nil, constant.NewError(
			-1,
			constant.ResponseErr,
			"响应缺少daily字段",
		)
	}

	dailyList, ok := daily.([]interface{})
	if !ok || len(dailyList) == 0 {
		return nil, constant.NewError(
			-1,
			constant.ResponseErr,
			"daily格式无效",
		)
	}

	return dailyList, nil
}

// QueryHourlyWeather 查询逐小时天气
func (h *APIHandler) QueryHourlyWeather(position *request.QueryWeatherLonLatRequest) ([]interface{}, error) {
	url := fmt.Sprintf("%s%s,%s", config.QueryHourWeatherAPI, position.Longitude, position.Latitude)
	response, err := h.doRequest(url)
	if err != nil {
		return nil, fmt.Errorf("逐小时天气查询失败: %w", err)
	}

	return extractHourlyWeather(response)
}

func extractHourlyWeather(response map[string]interface{}) ([]interface{}, error) {
	hourly, ok := response["hourly"]
	if !ok {
		return nil, constant.NewError(
			-1,
			constant.ResponseErr,
			"响应缺少hourly字段",
		)
	}

	hourlyList, ok := hourly.([]interface{})
	if !ok || len(hourlyList) == 0 {
		return nil, constant.NewError(
			-1,
			constant.ResponseErr,
			"hourly格式无效",
		)
	}

	return hourlyList, nil
}

// SearchPositionWithName 查询位置信息

func (h *APIHandler) SearchPosition(posInfo *request.QueryPositionRequest) ([]interface{}, error) {
	url := fmt.Sprintf("%s%s", config.UrbanSearchAPI, posInfo.Position)
	//url += "&lang=zh"
	response, err := h.doRequest(url)
	if err != nil {
		return nil, fmt.Errorf("位置查询失败: %w", err)
	}
	return extractPosInfo(response)
}

func extractPosInfo(response map[string]interface{}) ([]interface{}, error) {
	location, ok := response["location"]
	if !ok {
		return nil, constant.NewError(
			-1,
			constant.ResponseErr,
			"响应缺少location字段",
		)
	}

	locations, ok := location.([]interface{})
	if !ok || len(locations) == 0 {
		return nil, constant.NewError(
			-1,
			constant.ResponseErr,
			"location格式无效",
		)
	}
	return locations, nil
}
