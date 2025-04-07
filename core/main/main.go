package main

import (
	"core/config"
	"core/constant"
	"core/util"
	"gorm.io/gorm"
	"log"
)

func main() {
	// 初始化数据库连接
	util.InitDB()
	// 使用defer确保在函数结束时关闭数据库连接
	defer func(db *gorm.DB) {
		if db != nil {
			sqlDB, err := db.DB()
			if err != nil {
				panic(constant.DataBaseCannotBeCorrectlyClosed + err.Error())
			}
			err = sqlDB.Close()
			if err != nil {
				panic(constant.DataBaseCannotBeCorrectlyClosed + err.Error())
			}
		}
	}(config.DataBase)
	engine := util.InitGin()
	InitRouter(engine)
	err := engine.Run(config.RunPort)
	if err != nil {
		log.Fatal(err)
		//panic(err)
	}
}
