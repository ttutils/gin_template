package bootstrao

import (
	"gin_template/biz/model"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// 自动迁移表结构
	if err := db.AutoMigrate(
		&model.User{},
		&model.TenantInfo{},
	); err != nil {
		return err
	}

	err := InitData(db)
	if err != nil {
		return err
	}

	return nil
}
