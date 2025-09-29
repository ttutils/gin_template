package bootstrao

import (
	"gin_template/biz/model"
	"gin_template/utils"
	"gin_template/utils/config"

	"github.com/gookit/slog"
	"gorm.io/gorm"
)

func InitData(db *gorm.DB) error {
	// 插入初始化账号
	var count int64
	if err := db.Model(&model.User{}).Where("id = ?", 1).Count(&count).Error; err != nil {
		return err
	}

	// 如果不存在则创建
	if count == 0 {
		slog.Infof("%s 用户不存在，密码为:%s", config.Cfg.Admin.Username, config.Cfg.Admin.Password)
		adminUser := &model.User{
			Username: config.Cfg.Admin.Username,
			Password: utils.MD5(config.Cfg.Admin.Password),
		}
		if err := db.Create(adminUser).Error; err != nil {
			return err
		}
	}

	return nil
}
