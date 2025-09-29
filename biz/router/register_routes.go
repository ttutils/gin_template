package router

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	registerDiyRoutes(r)
	tenantRoutes(r)
	userRoutes(r)
}
