package router

import (
	hUser "gin_template/biz/handler/user"

	"github.com/gin-gonic/gin"
)

func userRoutes(r *gin.Engine) {
	userGroup := r.Group("/api")
	{
		userGroup.PUT("/user/add", hUser.CreateUser)
		userGroup.DELETE("/user/delete/:user_id", hUser.DeleteUser)
		userGroup.POST("/user/update/:user_id", hUser.UpdateUser)
		userGroup.POST("/user/change_passwd/:user_id", hUser.ChangePasswd)
		userGroup.POST("/user/login", hUser.UserLogin)
		userGroup.GET("/user/list", hUser.UserList)
		userGroup.GET("/user/info/:user_id", hUser.UserInfo)
	}
}
