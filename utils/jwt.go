package utils

import (
	"fmt"
	"gin_template/utils/config"
	"time"

	"github.com/gookit/slog"

	ginjwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtMiddleware *ginjwt.GinJWTMiddleware
	jwtSecret     = []byte(config.Cfg.Jwt.Secret)
)

// 初始化 JWT 中间件
func initJWT() error {
	var err error
	jwtMiddleware, err = ginjwt.New(&ginjwt.GinJWTMiddleware{
		Key:               jwtSecret,
		IdentityKey:       "userid",
		SendCookie:        false,
		SendAuthorization: false,
		TokenLookup:       "header: Authorization",
		TokenHeadName:     "Bearer",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(map[string]interface{}); ok {
				return jwt.MapClaims{
					"userid":   v["userid"],
					"username": v["username"],
					"iss":      config.Cfg.Server.Name,
				}
			}
			return jwt.MapClaims{}
		},
		// 添加登录验证函数，这是必需的
		Authenticator: func(c *gin.Context) (interface{}, error) {
			// 这里我们只是返回数据，实际使用时需要验证用户名密码
			return map[string]interface{}{
				"userid":   1,
				"username": "admin",
			}, nil
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

	// 暂存原超时配置
	originalTimeout := jwtMiddleware.Timeout
	defer func() { jwtMiddleware.Timeout = originalTimeout }()

	// 设置新的超时时间
	if len(expTime) > 0 {
		jwtMiddleware.Timeout = time.Second * time.Duration(expTime[0])
	} else {
		jwtMiddleware.Timeout = time.Hour * time.Duration(config.Cfg.Jwt.ExpireTime)
	}

	// 准备用户数据
	loginData := map[string]interface{}{
		"userid":   int(userid),
		"username": username,
	}

	// 直接调用 gin-jwt 的 PayloadFunc 来获取包含自定义字段的 claims
	// 这是确保 PayloadFunc 被调用的可靠方法
	claims := jwtMiddleware.PayloadFunc(loginData)

	// PayloadFunc 不会自动添加时间相关字段，需要手动添加
	claims["exp"] = time.Now().Add(jwtMiddleware.Timeout).Unix()
	claims["orig_iat"] = time.Now().Unix()

	// 使用标准的 jwt 库生成 token，但 claims 来自 gin-jwt 的 PayloadFunc
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
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

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("无法解析 claims")
	}

	// 验证 issuer
	if iss, ok := claims["iss"].(string); !ok || iss != config.Cfg.Server.Name {
		return nil, fmt.Errorf("issuer 不匹配")
	}

	// 验证过期时间
	if exp, ok := claims["exp"].(float64); !ok || int64(exp) < time.Now().Unix() {
		return nil, fmt.Errorf("token 已过期")
	}

	// userid 转换为 int
	if useridVal, ok := claims["userid"]; ok {
		if userid, ok := useridVal.(float64); ok {
			claims["userid"] = int(userid)
		}
	}

	// 验证 username
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

	// 保存 claims 到 gin.Context
	for k, v := range claims {
		c.Set(k, v)
	}
	c.Set("userid", claims["userid"])
	c.Set("username", claims["username"])

	// 检查是否管理员
	return IsAdmin(c)
}
