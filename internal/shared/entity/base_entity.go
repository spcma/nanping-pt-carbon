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
	CreateBy   int64         `gorm:"column:create_by;default:0" json:"createBy"`
	UpdateBy   int64         `gorm:"column:update_by;default:0" json:"updateBy"`
	DeleteBy   int64         `gorm:"column:delete_by;default:0" json:"deleteBy"`
	CreateTime timeutil.Time `gorm:"column:create_time;type:timestamp;default:timezone('Asia/Shanghai'::text, now())" json:"createTime"`
	UpdateTime timeutil.Time `gorm:"column:update_time;type:timestamp;default:timezone('Asia/Shanghai'::text, now())" json:"updateTime"`
	DeleteTime timeutil.Time `gorm:"column:delete_time;type:timestamp;default:'1970-01-01 00:00:00+08'" json:"deleteTime"`
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
