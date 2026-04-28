package domain

import (
	"app/internal/shared/entity"
	"context"
)

// ProjectMembersRepository 项目成员仓储接口（领域层定义）
type ProjectMembersRepository interface {
	Create(ctx context.Context, member *ProjectMembers) error
	Update(ctx context.Context, member *ProjectMembers) error
	Delete(ctx context.Context, id int64, userID int64) error

	FindByID(ctx context.Context, id int64) (*ProjectMembers, error)
	FindByProjectID(ctx context.Context, projectID int64) ([]*ProjectMembers, error)
	FindByUserID(ctx context.Context, userID int64) ([]*ProjectMembers, error)
	FindByProjectAndUser(ctx context.Context, projectID, userID int64) (*ProjectMembers, error)
	FindPage(ctx context.Context, query *ProjectMembersPageQuery) (*entity.PaginationResult[*ProjectMembers], error)
}

// ProjectMembersPageQuery 项目成员分页查询对象
type ProjectMembersPageQuery struct {
	entity.PaginationQuery
	ProjectId int64  `json:"projectId" form:"projectId"`
	UserId    int64  `json:"userId" form:"userId"`
	Role      string `json:"role" form:"role"`
	SortBy    string `json:"sortBy" form:"sortBy"`
	SortOrder string `json:"sortOrder" form:"sortOrder"` // "asc" or "desc"
}
