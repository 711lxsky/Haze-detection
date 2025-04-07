package model

import (
	"time"
)

// QueryRecord 表示查询记录的结构体
type QueryRecord struct {
	ID          int64     `db:"id" json:"id"`
	City        string    `db:"city" json:"city"`
	Longitude   string    `db:"longitude" json:"longitude"`
	Latitude    string    `db:"latitude" json:"latitude"`
	WeatherInfo string    `db:"weather_info" json:"weather_info"`
	CreateTime  time.Time `db:"create_time" json:"create_time"`
	Type        int8      `db:"type" json:"type"`
	Deleted     int8      `db:"deleted" json:"deleted"`
}
