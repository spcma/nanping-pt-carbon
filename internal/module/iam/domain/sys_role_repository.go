package domain

import (
	"context"
)

// SysRolePageQuery system role page query object
type SysRolePageQuery struct {
	PageNum   int64  `json:"pageNum" binding:"required,min=1"`
	PageSize  int64  `json:"pageSize" binding:"required,min=1,max=100"`
	Name      string `json:"name"`   // 角色名模糊匹配
	Code      string `json:"code"`   // 角色编码精确匹配
	Status    string `json:"status"` // 状态
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder" binding:"oneof=asc desc"` // "asc" or "desc"
}

// SysRoleRepository role repository interface
type SysRoleRepository interface {
	Create(ctx context.Context, role *SysRole) error
	Update(ctx context.Context, role *SysRole) error
	Delete(ctx context.Context, id int64) error
	FindByID(ctx context.Context, id int64) (*SysRole, error)
	FindByCode(ctx context.Context, code string) (*SysRole, error)
	FindPage(ctx context.Context, query SysRolePageQuery) ([]*SysRole, int64, error)
	FindListByCodes(ctx context.Context, codes []string) ([]*SysRole, error)
}
