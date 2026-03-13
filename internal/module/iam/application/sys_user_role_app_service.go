package application

import (
	"app/internal/module/iam/domain"
	"context"
)

// AssignUserRoleCommand assign user role command
type AssignUserRoleCommand struct {
	UserID   int64 `json:"userId"`
	RoleID   int64 `json:"roleId"`
	UserIDOp int64 `json:"userIdOp"`
}

// CancelUserRoleCommand cancel user role command
type CancelUserRoleCommand struct {
	UserID   int64 `json:"userId"`
	RoleID   int64 `json:"roleId"`
	UserIDOp int64 `json:"userIdOp"`
}

// SysUserRoleAppService system user role application service
type SysUserRoleAppService struct {
	repo domain.SysUserRoleRepository
}

// NewSysUserRoleAppService creates system user role application service
func NewSysUserRoleAppService(repo domain.SysUserRoleRepository) *SysUserRoleAppService {
	return &SysUserRoleAppService{repo: repo}
}

// AssignUserRole assigns a role to a user
func (s *SysUserRoleAppService) AssignUserRole(ctx context.Context, cmd AssignUserRoleCommand) error {
	userRole, err := domain.NewSysUserRole(cmd.UserID, cmd.RoleID, cmd.UserIDOp)
	if err != nil {
		return err
	}
	return s.repo.Create(ctx, userRole)
}

// CancelUserRole cancels a role from a user
func (s *SysUserRoleAppService) CancelUserRole(ctx context.Context, cmd CancelUserRoleCommand) error {
	return s.repo.Delete(ctx, cmd.UserID, cmd.RoleID)
}

// GetUserRoles gets roles of a user
func (s *SysUserRoleAppService) GetUserRoles(ctx context.Context, userId int64) ([]*domain.SysUserRole, error) {
	return s.repo.FindByUserID(ctx, userId)
}

// GetRoleUsers gets users of a role
func (s *SysUserRoleAppService) GetRoleUsers(ctx context.Context, roleId int64) ([]*domain.SysUserRole, error) {
	return s.repo.FindByRoleID(ctx, roleId)
}
