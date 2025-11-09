package router

import (
	hTenant "gin_template/biz/handler/tenant"
	"gin_template/biz/mw"

	"github.com/gin-gonic/gin"
)

func tenantRoutes(r *gin.RouterGroup) {
	tenantGroup := r.Group("/tenant")
	tenantGroup.Use(mw.JWTAuthMiddleware())
	{
		tenantGroup.PUT("/add", hTenant.CreateTenant)
		tenantGroup.DELETE("/delete/:id", hTenant.DeleteTenant)
		tenantGroup.GET("/list", hTenant.TenantList)
	}
}
