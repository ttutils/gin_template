package utils

import (
	"fmt"
	"gin_template/utils/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gookit/slog"
)

var (
	// JWT 配置
	jwtConfig = struct {
		Secret        []byte
		IdentityKey   string
		TokenHeadName string
		Issuer        string
	}{
		Secret:        []byte(config.Cfg.Jwt.Secret),
		IdentityKey:   "userid",
		TokenHeadName: "Bearer",
	}
)

// GenerateToken 生成 JWT Token
func GenerateToken(userid uint, username string, expTime ...int) (string, error) {
	// 设置过期时间
	var expireTime time.Duration
	if len(expTime) > 0 {
		expireTime = time.Minute * time.Duration(expTime[0])
	} else {
		expireTime = time.Hour * time.Duration(config.Cfg.Jwt.ExpireTime)
	}

	// 创建 claims
	claims := jwt.MapClaims{
		"userid":   int(userid),
		"username": username,
		"iss":      config.Cfg.Server.Name,
		"exp":      time.Now().Add(expireTime).Unix(),
		"orig_iat": time.Now().Unix(),
	}

	// 生成 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtConfig.Secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析并验证 JWT Token
func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	// 解析 token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtConfig.Secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("token 解析失败: %v", err)
	}

	// 验证 token 是否有效
	if !token.Valid {
		return nil, fmt.Errorf("token 无效")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
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
		return fmt.Errorf("token 无效: %v", err)
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
