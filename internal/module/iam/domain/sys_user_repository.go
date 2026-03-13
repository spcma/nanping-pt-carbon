package domain

import (
	"context"
)

// SysUserPageQuery system user page query object
type SysUserPageQuery struct {
	PageNum   int64  `json:"pageNum" binding:"required,min=1"`
	PageSize  int64  `json:"pageSize" binding:"required,min=1,max=100"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	UserType  string `json:"userType"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder" binding:"oneof=asc desc"` // "asc" or "desc"
}

// SysUserRepository
type SysUserRepository interface {
	Create(ctx context.Context, user *SysUser) error
	Update(ctx context.Context, user *SysUser) error
	Delete(ctx context.Context, id int64, userID int64) error // 逻辑删除
	FindByID(ctx context.Context, id int64) (*SysUser, error)
	FindByUsername(ctx context.Context, username string) (*SysUser, error)
	FindPage(ctx context.Context, query *SysUserPageQuery) ([]*SysUser, int64, error)
}
