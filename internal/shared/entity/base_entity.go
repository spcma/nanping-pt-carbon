package entity

import (
	"app/internal/shared/timeutil"
)

// 公共字段名称常量，便于统一管理和维护
const (
	FieldId         = "id"
	FieldCreateBy   = "create_by"
	FieldUpdateBy   = "update_by"
	FieldDeleteBy   = "delete_by"
	FieldCreateTime = "create_time"
	FieldUpdateTime = "update_time"
	FieldDeleteTime = "delete_time"
)

// BaseEntity 通用基础实体，包含所有实体共有的审计字段
type BaseEntity struct {
	Id         int64         `gorm:"primaryKey;default:0;autoIncrement:false" json:"id"`
	CreateBy   int64         `gorm:"column:create_by;default:0" json:"create_by"`
	UpdateBy   int64         `gorm:"column:update_by;default:0" json:"update_by"`
	DeleteBy   int64         `gorm:"column:delete_by;default:0" json:"delete_by"`
	CreateTime timeutil.Time `gorm:"column:create_time;type:timestamp;default:timezone('Asia/Shanghai'::text, now())" json:"create_time"`
	UpdateTime timeutil.Time `gorm:"column:update_time;type:timestamp;default:timezone('Asia/Shanghai'::text, now())" json:"update_time"`
	DeleteTime timeutil.Time `gorm:"column:delete_time;type:timestamp;default:'1970-01-01 00:00:00+08'" json:"delete_time"`
}

// Delete 逻辑删除
func (b *BaseEntity) Delete(userID int64) error {
	b.DeleteBy = userID
	b.DeleteTime = timeutil.Now()

	return nil
}

// IsDeleted 判断是否逻辑删除
func (b *BaseEntity) IsDeleted() bool {
	return b.DeleteBy > 0
}
