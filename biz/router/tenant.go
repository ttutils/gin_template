package router

import (
	hTenant "gin_template/biz/handler/tenant"

	"github.com/gin-gonic/gin"
)

func tenantRoutes(r *gin.Engine) {
	tenantGroup := r.Group("/api")
	{
		tenantGroup.PUT("/tenant/add", hTenant.CreateTenant)
		tenantGroup.DELETE("/tenant/delete/:id", hTenant.DeleteTenant)
		tenantGroup.GET("/tenant/list", hTenant.TenantList)
	}
}
