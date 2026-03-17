package db

import (
	"errors"
	"gorm.io/gorm"
)

// ErrFilter 过滤可接受的错误
// 将 gorm.ErrRecordNotFound 等预期内的错误转换为 nil
// 适用于单条记录查询场景，避免返回不必要的错误
func ErrFilter(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return err
}
