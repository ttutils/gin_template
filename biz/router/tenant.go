package router

import (
	hTenant "gin_template/biz/handler/tenant"

	"github.com/gin-gonic/gin"
)

func tenantRoutes(r *gin.Engine) {
	r.PUT("/api/tenant/add", hTenant.CreateTenant)
	r.DELETE("/api/tenant/delete/:id", hTenant.DeleteTenant)
	r.GET("/api/tenant/list", hTenant.TenantList)
}
