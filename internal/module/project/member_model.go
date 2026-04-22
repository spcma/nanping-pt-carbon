package project

import (
	"app/internal/shared/entity"
	idgen "app/internal/shared/idgen"
	"app/internal/shared/timeutil"
	"errors"
)

// ProjectMemberRole 项目成员角色
type ProjectMemberRole string

const (
	ProjectMemberRoleOwner  ProjectMemberRole = "owner"  // 所有者
	ProjectMemberRoleAdmin  ProjectMemberRole = "admin"  // 管理员
	ProjectMemberRoleMember ProjectMemberRole = "member" // 普通成员
)

// ProjectMembers 项目成员聚合根
type ProjectMembers struct {
	entity.BaseEntity
	ProjectId int64             `json:"project_id"`
	UserId    int64             `json:"user_id"`
	Role      ProjectMemberRole `json:"role" gorm:"column:role"` // owner admin member
}

// TableName 表名
func (*ProjectMembers) TableName() string {
	return "project_members"
}

// NewProjectMembers 添加项目新成员
func NewProjectMembers(projectId, userId, createUser int64, role string) (*ProjectMembers, error) {
	if role == "" {
		role = string(ProjectMemberRoleMember)
	}

	// 验证角色类型
	switch ProjectMemberRole(role) {
	case ProjectMemberRoleOwner, ProjectMemberRoleAdmin, ProjectMemberRoleMember:
		// 有效角色
	default:
		return nil, errors.New("invalid role type")
	}

	projectMember := &ProjectMembers{
		BaseEntity: entity.BaseEntity{
			Id:         idgen.NumID(),
			CreateBy:   createUser,
			CreateTime: timeutil.Now(),
		},
		ProjectId: projectId,
		UserId:    userId,
		Role:      ProjectMemberRole(role),
	}

	return projectMember, nil
}

// UpdateRole 更新成员角色
func (pm *ProjectMembers) UpdateRole(role string, userID int64) error {
	switch ProjectMemberRole(role) {
	case ProjectMemberRoleOwner, ProjectMemberRoleAdmin, ProjectMemberRoleMember:
		pm.Role = ProjectMemberRole(role)
		pm.UpdateBy = userID
		pm.UpdateTime = timeutil.Now()
	default:
		return errors.New("invalid role type")
	}
	return nil
}

// Delete 逻辑删除成员
func (pm *ProjectMembers) Delete(userID int64) error {
	pm.DeleteBy = userID
	pm.DeleteTime = timeutil.Now()
	return nil
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
