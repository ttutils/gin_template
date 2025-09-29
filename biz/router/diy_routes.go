package router

import (
	"gin_template/biz/handler"

	"github.com/gin-gonic/gin"
)

func registerDiyRoutes(r *gin.Engine) {
	DiyGroup := r.Group("/api")
	{
		DiyGroup.POST("/ping", handler.Ping)
		DiyGroup.GET("/server_info", handler.ServerInfo)
		DiyGroup.GET("/is_demo", handler.GetDemo)
		DiyGroup.GET("/metrics", handler.Metrics)
	}
}
