package migration

import (
	"app/internal/module/iam/domain"

	"gorm.io/gorm"
)

// MigrateTables 执行数据库迁移，创建所有表
func MigrateTables(db *gorm.DB) error {
	// 迁移用户表
	err := db.AutoMigrate(&domain.Users{})
	if err != nil {
		return err
	}

	// 迁移角色表
	err = db.AutoMigrate(&domain.SysRole{})
	if err != nil {
		return err
	}

	// 迁移用户角色关联表
	err = db.AutoMigrate(&domain.SysUserRole{})
	if err != nil {
		return err
	}

	return nil
}
