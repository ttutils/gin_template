package router

import (
	hUser "gin_template/biz/handler/user"

	"github.com/gin-gonic/gin"
)

func userRoutes(r *gin.Engine) {
	r.PUT("/api/user/add", hUser.CreateUser)
	r.DELETE("/api/user/delete/:user_id", hUser.DeleteUser)
	r.POST("/api/user/update/:user_id", hUser.UpdateUser)
	r.POST("/api/user/change_passwd/:user_id", hUser.ChangePasswd)
	r.POST("/api/user/login", hUser.UserLogin)
	r.POST("/nacos/v1/auth/login", hUser.NacosUserLogin)
	r.GET("/api/user/list", hUser.UserList)
	r.GET("/api/user/info/:user_id", hUser.UserInfo)
}
