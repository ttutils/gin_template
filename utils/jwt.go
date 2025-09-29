package utils

import (
	"fmt"
	"gin_template/utils/config"
	"time"

	"github.com/gookit/slog"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var (
	jwtMiddleware *jwt.GinJWTMiddleware
	jwtSecret     = []byte(config.Cfg.Jwt.Secret)
)

// 初始化 JWT 中间件
func initJWT() error {
	var err error
	jwtMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Key:               jwtSecret,
		IdentityKey:       "userid",
		SendCookie:        false,
		SendAuthorization: false,
		TokenLookup:       "header: Authorization",
		TokenHeadName:     "Bearer",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(map[string]interface{}); ok {
				return jwt.MapClaims{
					"userid":   v["userid"], // 确保生成时为int类型
					"username": v["username"],
					"iss":      config.Cfg.Server.Name,
				}
			}
			return jwt.MapClaims{}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if claims, ok := data.(jwt.MapClaims); ok {
				// 确保 userid 转换为 int 类型
				if userid, ok := claims["userid"].(float64); ok {
					// 设置到gin上下文中
					c.Set("userid", int(userid))
				}
				if username, ok := claims["username"].(string); ok {
					c.Set("username", username)
				}
				return true
			}
			return false
		},
	})
	return err
}

// GenerateToken 生成 JWT Token
func GenerateToken(userid uint, username string, expTime ...int) (string, error) {
	if jwtMiddleware == nil {
		if err := initJWT(); err != nil {
			return "", err
		}
	}

	// 强制覆盖中间件的 Timeout 设置
	originalTimeout := jwtMiddleware.Timeout
	defer func() { jwtMiddleware.Timeout = originalTimeout }() // 恢复原配置

	// 计算过期时间
	if len(expTime) > 0 {
		jwtMiddleware.Timeout = time.Hour * time.Duration(expTime[0])
	} else {
		jwtMiddleware.Timeout = time.Hour * time.Duration(config.Cfg.Jwt.ExpireTime)
	}

	// 生成 claims（不再手动设置 exp）
	claims := map[string]interface{}{
		"userid":   int(userid),
		"username": username,
	}

	token, _, err := jwtMiddleware.TokenGenerator(claims)
	return token, err
}

// ParseToken 解析并验证 JWT Token
func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	if jwtMiddleware == nil {
		if err := initJWT(); err != nil {
			return nil, err
		}
	}

	parsedToken, err := jwtMiddleware.ParseTokenString(tokenStr)
	if err != nil {
		return nil, fmt.Errorf("token 解析失败: %v", err)
	}

	claims := jwt.ExtractClaimsFromToken(parsedToken)

	// 验证 issuer
	if iss, ok := claims["iss"].(string); !ok || iss != config.Cfg.Server.Name {
		return nil, fmt.Errorf("issuer 不匹配")
	}

	// 验证过期时间
	if exp, ok := claims["exp"].(float64); !ok || int64(exp) < time.Now().Unix() {
		return nil, fmt.Errorf("token 已过期")
	}

	// 验证userid并转换为int
	useridVal, ok := claims["userid"]
	if !ok {
		return nil, fmt.Errorf("token 缺少 userid")
	}
	userid, ok := useridVal.(float64)
	if !ok {
		return nil, fmt.Errorf("userid 类型错误")
	}
	claims["userid"] = int(userid) // 转换为int类型

	// 验证username
	if _, ok := claims["username"].(string); !ok {
		return nil, fmt.Errorf("token 缺少 username")
	}

	return claims, nil
}

// GetUsernameFromContext 从上下文中提取用户名
func GetUsernameFromContext(c *gin.Context) (string, error) {
	usernameVal, exists := c.Get("username")
	if !exists {
		return "", fmt.Errorf("未找到用户名")
	}

	username, ok := usernameVal.(string)
	if !ok {
		return "", fmt.Errorf("用户名类型错误")
	}

	return username, nil
}

// GetUseridFromContext 从上下文中提取用户ID
func GetUseridFromContext(c *gin.Context) (int, error) {
	useridVal, exists := c.Get("userid")
	if !exists {
		return 0, fmt.Errorf("未找到用户ID")
	}

	userid, ok := useridVal.(int)
	if !ok {
		return 0, fmt.Errorf("userid 类型错误")
	}

	return userid, nil
}

// IsAdmin 判断是否为管理员
func IsAdmin(c *gin.Context) error {
	userid, err := GetUseridFromContext(c)
	if err != nil {
		return err
	}

	if userid != 1 {
		slog.Infof("当前用户ID: %d", userid)
		return fmt.Errorf("不是管理员，没有权限")
	}

	return nil
}

func ValidateToken(c *gin.Context, token string) error {
	// 验证 token
	claims, err := ParseToken(token)
	if err != nil {
		return fmt.Errorf("token 无效")
	}

	// 将 claims 保存到上下文
	for k, v := range claims {
		c.Set(k, v)
	}
	c.Set("userid", claims["userid"])
	c.Set("username", claims["username"])

	// 检查是否为管理员
	return IsAdmin(c)
}
