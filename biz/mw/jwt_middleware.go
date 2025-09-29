package mw

import (
	"gin_template/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 鉴权中间件
func JWTAuthMiddleware(excludedPaths []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求路径并转换为字符串
		path := c.Request.URL.Path

		skipPrefixes := []string{"/api"}

		// 需要鉴权的路径
		for _, prefix := range skipPrefixes {
			if !strings.HasPrefix(path, prefix) {
				c.Next()
				return
			}
		}

		// 如果路径是 swagger 文档，则跳过鉴权
		if strings.HasPrefix(path, "/api/swagger") {
			c.Next()
			return
		}

		// 如果路径在排除列表中，则跳过鉴权
		for _, excludedPath := range excludedPaths {
			if path == excludedPath {
				c.Next() // 跳过中间件，直接处理请求
				return
			}
		}

		// 获取 Authorization Header
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code": http.StatusUnauthorized,
				"msg":  "缺少token",
			})
			c.Abort() // 终止后续处理
			return
		}

		// 提取token（去除Bearer前缀）
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code": http.StatusUnauthorized,
				"msg":  "token格式错误",
			})
			c.Abort() // 终止后续处理
			return
		}

		// 验证 token
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code": http.StatusUnauthorized,
				"msg":  err.Error(),
			})
			c.Abort() // 终止后续处理
			return
		}

		// 将 claims 保存到上下文
		for k, v := range claims {
			c.Set(k, v)
		}
		c.Set("userid", claims["userid"])
		c.Set("username", claims["username"])

		// 如果验证通过，继续处理请求
		c.Next()
	}
}
