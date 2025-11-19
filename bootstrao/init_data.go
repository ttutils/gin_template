package bootstrao

import (
	"gin_template/biz/model"
	"gin_template/utils"
	"gin_template/utils/config"

	"github.com/gookit/slog"
	"gorm.io/gorm"
)

func InitData(db *gorm.DB) error {
	adminUser := &model.User{
		Username: config.Cfg.Admin.Username,
		Password: utils.MD5(config.Cfg.Admin.Password),
		Enable:   true,
	}

	result := db.Where(model.User{Username: config.Cfg.Admin.Username}).FirstOrCreate(adminUser)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
		slog.Infof("创建管理员用户成功，用户名: %s, 密码: %s", config.Cfg.Admin.Username, config.Cfg.Admin.Password)
	} else {
		slog.Infof("管理员用户已存在: %s", config.Cfg.Admin.Username)
	}

	return nil
}
