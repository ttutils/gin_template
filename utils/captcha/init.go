package captcha

import (
	"gin_template/utils/config"
	"time"

	"github.com/mojocn/base64Captcha"
)

// Driver 验证码驱动 - 延迟初始化
var Driver *base64Captcha.DriverDigit
var Store base64Captcha.Store

// Init 初始化验证码配置
func Init() {
	Driver = base64Captcha.NewDriverDigit(
		60,                          // 高度
		240,                         // 宽度
		config.Cfg.Captcha.Length,   // 长度
		config.Cfg.Captcha.MaxSkew,  // 最大倾斜度
		config.Cfg.Captcha.DotCount, // 点数
	)
	// 验证码存储 - 设置5分钟过期
	Store = base64Captcha.NewMemoryStore(10240, time.Duration(config.Cfg.Server.CaptchaExpireTime)*time.Minute)
}
