package router

import (
	hTenant "gin_template/biz/handler/tenant"

	"github.com/gin-gonic/gin"
)

func tenantRoutes(r *gin.RouterGroup) {
	tenantGroup := r.Group("/tenant")
	{
		tenantGroup.PUT("/add", hTenant.CreateTenant)
		tenantGroup.DELETE("/delete/:id", hTenant.DeleteTenant)
		tenantGroup.GET("/list", hTenant.TenantList)
	}
}
