package router

import (
	hUser "gin_template/biz/handler/user"

	"github.com/gin-gonic/gin"
)

func userRoutes(apiGroup *gin.RouterGroup) {
	userGroup := apiGroup.Group("/user")
	{
		userGroup.PUT("/add", hUser.CreateUser)
		userGroup.DELETE("/delete/:user_id", hUser.DeleteUser)
		userGroup.POST("/update/:user_id", hUser.UpdateUser)
		userGroup.POST("/change_passwd/:user_id", hUser.ChangePasswd)
		userGroup.POST("/login", hUser.UserLogin)
		userGroup.GET("/list", hUser.UserList)
		userGroup.GET("/info/:user_id", hUser.UserInfo)
	}
}
