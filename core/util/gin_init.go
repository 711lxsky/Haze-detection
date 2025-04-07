package util

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

func InitGin() *gin.Engine {
	// 创建
	engine := gin.Default()
	// 跨域调用设置
	crossConfig(engine)
	return engine
}

func crossConfig(r *gin.Engine) {
	// 添加 CORS 中间件， 允许跨域请求访问
	r.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}
