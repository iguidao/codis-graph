package http

import (
	//"net/http"

	v1 "codisgraph/src/http/v1"

	"github.com/gin-gonic/gin"
)

// NewServer return a configured http server of gin
func NewServer() *gin.Engine {

	// 存储日志文件代码
	// logpath := cfg.Get_Local("logpath")
	// gin.DisableConsoleColor()
	// f, _ := os.Create(logpath)
	// gin.DefaultWriter = io.MultiWriter(f)
	r := gin.Default()

	// 验证相关
	api := r.Group("/api/v1")
	{
		api.GET("/health", v1.HealthCheck)
		api.GET("/cookie", v1.Cookie)
	}
	return r
}
