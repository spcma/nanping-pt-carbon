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

// UserRoleAppService system user role application service
type UserRoleAppService struct {
	repo UserRoleRepository
}

// NewUserRoleAppService creates system user role application service
func NewUserRoleAppService(repo UserRoleRepository) *UserRoleAppService {
	return &UserRoleAppService{repo: repo}
}

// AssignUserRole assigns a role to a user
func (s *UserRoleAppService) AssignUserRole(ctx context.Context, cmd AssignUserRoleCommand) error {
	userRole, err := domain.NewUserRole(cmd.UserID, cmd.RoleID, cmd.UserIDOp)
	if err != nil {
		return err
	}
	return s.repo.Create(ctx, userRole)
}

// CancelUserRole cancels a role from a user
func (s *UserRoleAppService) CancelUserRole(ctx context.Context, cmd CancelUserRoleCommand) error {
	return s.repo.DeleteByUserRoleID(ctx, cmd.UserID, cmd.RoleID)
}

// GetUserRoles gets roles of a user
func (s *UserRoleAppService) GetUserRoles(ctx context.Context, userId int64) ([]*domain.UserRole, error) {
	return s.repo.FindByUserID(ctx, userId)
}

// GetRoleUsers gets users of a role
func (s *UserRoleAppService) GetRoleUsers(ctx context.Context, roleId int64) ([]*domain.UserRole, error) {
	return s.repo.FindByRoleID(ctx, roleId)
}
