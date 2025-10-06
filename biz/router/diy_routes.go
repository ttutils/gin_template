package router

import (
	"gin_template/biz/handler"

	"github.com/gin-gonic/gin"
)

func registerDiyRoutes(apiGroup *gin.RouterGroup) {
	apiGroup.GET("/ping", handler.Ping)
	apiGroup.GET("/server_info", handler.ServerInfo)
	apiGroup.GET("/is_demo", handler.GetDemo)
	apiGroup.GET("/metrics", handler.Metrics)
}
