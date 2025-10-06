package router

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	apiGroup := r.Group("/api")
	diyRoutes(apiGroup)
	tenantRoutes(apiGroup)
	userRoutes(apiGroup)
}
