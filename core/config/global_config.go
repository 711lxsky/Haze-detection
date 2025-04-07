package config

import "gorm.io/gorm"

var DataBase *gorm.DB

var (
	RunPort = ":8248"

	DatabaseName          = "haze_detection"
	DatabaseUser          = "root"
	DatabasePassword      = "zyy_lhx_yjr"
	DatabaseHost          = "127.0.0.1"
	DatabasePort          = "3308"
	DatabaseConnectParams = "charset=utf8mb4&parseTime=True&loc=Local"
)

var (
	ResponseMessage = "massage"
	ResponseReason  = "reason"
	ResponseData    = "data"
	Success         = "success"
)
