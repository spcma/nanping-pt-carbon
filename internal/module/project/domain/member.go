package domain

import (
	"app/internal/shared/entity"
	idgen "app/internal/shared/idgen"
	"app/internal/shared/timeutil"
	"errors"
)

// ProjectMemberRole 项目成员角色（值对象）
type ProjectMemberRole string

const (
	ProjectMemberRoleOwner  ProjectMemberRole = "owner"  // 所有者
	ProjectMemberRoleAdmin  ProjectMemberRole = "admin"  // 管理员
	ProjectMemberRoleMember ProjectMemberRole = "member" // 普通成员
)

// IsValid 验证角色是否有效
func (r ProjectMemberRole) IsValid() bool {
	switch r {
	case ProjectMemberRoleOwner, ProjectMemberRoleAdmin, ProjectMemberRoleMember:
		return true
	default:
		return false
	}
}

// ProjectMembers 项目成员聚合根
type ProjectMembers struct {
	entity.BaseEntity
	ProjectId int64             `json:"project_id"`
	UserId    int64             `json:"user_id"`
	Role      ProjectMemberRole `json:"role" gorm:"column:role"`
}

// TableName 表名
func (*ProjectMembers) TableName() string {
	return "project_members"
}

// NewProjectMembers 创建新项目成员（工厂方法）
func NewProjectMembers(projectId, userId, createUser int64, role string) (*ProjectMembers, error) {
	// 默认角色
	if role == "" {
		role = string(ProjectMemberRoleMember)
	}

	// 验证角色类型
	memberRole := ProjectMemberRole(role)
	if !memberRole.IsValid() {
		return nil, errors.New("无效的成员角色")
	}

	projectMember := &ProjectMembers{
		BaseEntity: entity.BaseEntity{
			Id:         idgen.NumID(),
			CreateBy:   createUser,
			CreateTime: timeutil.Now(),
		},
		ProjectId: projectId,
		UserId:    userId,
		Role:      memberRole,
	}

	return projectMember, nil
}

// UpdateRole 更新成员角色（领域行为）
func (pm *ProjectMembers) UpdateRole(role string, userID int64) error {
	memberRole := ProjectMemberRole(role)

	// 验证角色有效性
	if !memberRole.IsValid() {
		return errors.New("无效的成员角色")
	}

	// 执行更新
	pm.Role = memberRole
	pm.UpdateBy = userID
	pm.UpdateTime = timeutil.Now()

	return nil
}

// Delete 逻辑删除成员（领域行为）
func (pm *ProjectMembers) Delete(userID int64) error {
	// 领域规则：所有者不能被删除，只能转移所有权
	if pm.Role == ProjectMemberRoleOwner {
		return errors.New("不能删除项目所有者，请先转移所有权")
	}

	pm.DeleteBy = userID
	pm.DeleteTime = timeutil.Now()
	return nil
}

// IsOwner 检查是否为所有者
func (pm *ProjectMembers) IsOwner() bool {
	return pm.Role == ProjectMemberRoleOwner
}

// IsAdmin 检查是否为管理员
func (pm *ProjectMembers) IsAdmin() bool {
	return pm.Role == ProjectMemberRoleAdmin
}

// HasPermission 检查是否有指定权限
func (pm *ProjectMembers) HasPermission(requiredRole ProjectMemberRole) bool {
	// 权限等级：owner > admin > member
	roleLevel := map[ProjectMemberRole]int{
		ProjectMemberRoleMember: 1,
		ProjectMemberRoleAdmin:  2,
		ProjectMemberRoleOwner:  3,
	}

	return roleLevel[pm.Role] >= roleLevel[requiredRole]
}
