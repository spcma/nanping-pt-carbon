package application

import (
	"app/internal/module/iam/domain"
	"app/internal/shared/entity"
	"context"
)

// CreateSysRoleCommand create system role command
type CreateSysRoleCommand struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	UserID      int64  `json:"userId"`
}

// UpdateSysRoleCommand update system role command
type UpdateSysRoleCommand struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      int64  `json:"userId"`
}

// ChangeRoleStatusCommand change role status command
type ChangeRoleStatusCommand struct {
	ID     int64             `json:"id"`
	Status domain.RoleStatus `json:"status"`
	UserID int64             `json:"userId"`
}

// RoleAppService system role application service
type RoleAppService struct {
	repo RoleRepo
}

// NewRoleAppService creates system role application service
func NewRoleAppService(repo RoleRepo) *RoleAppService {
	return &RoleAppService{repo: repo}
}

// CreateRole creates a system role
func (s *RoleAppService) CreateRole(ctx context.Context, cmd CreateSysRoleCommand) (int64, error) {
	role, err := domain.NewRole(cmd.Name, cmd.Code, cmd.Description, cmd.UserID)
	if err != nil {
		return 0, err
	}
	err = s.repo.Create(ctx, role)
	if err != nil {
		return 0, err
	}
	return role.Id, nil
}

// UpdateRole updates a system role
func (s *RoleAppService) UpdateRole(ctx context.Context, cmd UpdateSysRoleCommand) error {
	role, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	return role.UpdateInfo(cmd.Name, cmd.Description, cmd.UserID)
}

// DeleteRole deletes a system role
func (s *RoleAppService) DeleteRole(ctx context.Context, id, uid int64) error {
	return s.repo.Delete(ctx, id, uid)
}

// GetRoleByID gets system role by ID
func (s *RoleAppService) GetRoleByID(ctx context.Context, id int64) (*domain.Role, error) {
	return s.repo.FindByID(ctx, id)
}

// GetRoleByCode gets system role by code
func (s *RoleAppService) GetRoleByCode(ctx context.Context, code string) (*domain.Role, error) {
	return s.repo.FindByCode(ctx, code)
}

// GetRolePage queries system roles with pagination
func (s *RoleAppService) GetRolePage(ctx context.Context, query *domain.SysRolePageQuery) (*entity.PaginationResult[*domain.Role], error) {
	result, err := s.repo.FindPage(ctx, query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ChangeRoleStatus changes role status
func (s *RoleAppService) ChangeRoleStatus(ctx context.Context, cmd ChangeRoleStatusCommand) error {
	role, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	return role.ChangeStatus(cmd.Status, cmd.UserID)
}
