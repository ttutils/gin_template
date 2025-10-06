package router

import (
	"gin_template/biz/handler"

	"github.com/gin-gonic/gin"
)

func registerDiyRoutes(r *gin.Engine) {
	diyGroup := r.Group("/api")
	{
		diyGroup.GET("/ping", handler.Ping)
		diyGroup.GET("/server_info", handler.ServerInfo)
		diyGroup.GET("/is_demo", handler.GetDemo)
		diyGroup.GET("/metrics", handler.Metrics)
	}
}
