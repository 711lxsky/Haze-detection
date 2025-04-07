package util

import (
	"core/config"
	"core/constant"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB 初始化数据库连接，并返回数据库实例
// 该函数使用了defer来确保在函数结束时关闭数据库连接
// 使用AutoMigrate方法来确保数据库模式与模型匹配
func InitDB() {
	// 打开数据库连接
	var err error
	// 拼接MySQL连接配置
	dbConnect := config.DatabaseUser + ":" + config.DatabasePassword + "@tcp(" + config.DatabaseHost + ":" + config.DatabasePort + ")/" + config.DatabaseName + "?" + config.DatabaseConnectParams
	config.DataBase, err = gorm.Open(mysql.Open(dbConnect), &gorm.Config{})
	if err != nil {
		panic(constant.CannotConnectDB + err.Error())
	}
	// 自动迁移模式， 保持更新到最新
	// 仅创建表， 缺少列和索引， 不会改变现有列的类型或删除未使用的列以保护数据
	//if err := config.DataBase.AutoMigrate(
	//	&model.QueryRecord{},
	//); err != nil {
	//	panic(err.Error())
	//}
}
