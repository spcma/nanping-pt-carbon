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

// SysRoleAppService system role application service
type SysRoleAppService struct {
	repo RoleRepo
}

// NewSysRoleAppService creates system role application service
func NewSysRoleAppService(repo RoleRepo) *SysRoleAppService {
	return &SysRoleAppService{repo: repo}
}

// CreateSysRole creates a system role
func (s *SysRoleAppService) CreateSysRole(ctx context.Context, cmd CreateSysRoleCommand) (int64, error) {
	role, err := domain.NewSysRole(cmd.Name, cmd.Code, cmd.Description, cmd.UserID)
	if err != nil {
		return 0, err
	}
	err = s.repo.Create(ctx, role)
	if err != nil {
		return 0, err
	}
	return role.Id, nil
}

// UpdateSysRole updates a system role
func (s *SysRoleAppService) UpdateSysRole(ctx context.Context, cmd UpdateSysRoleCommand) error {
	role, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	return role.UpdateInfo(cmd.Name, cmd.Description, cmd.UserID)
}

// DeleteSysRole deletes a system role
func (s *SysRoleAppService) DeleteSysRole(ctx context.Context, id, uid int64) error {
	return s.repo.Delete(ctx, id, uid)
}

// GetSysRoleByID gets system role by ID
func (s *SysRoleAppService) GetSysRoleByID(ctx context.Context, id int64) (*domain.SysRole, error) {
	return s.repo.FindByID(ctx, id)
}

// GetSysRoleByCode gets system role by code
func (s *SysRoleAppService) GetSysRoleByCode(ctx context.Context, code string) (*domain.SysRole, error) {
	return s.repo.FindByCode(ctx, code)
}

// GetSysRolePage queries system roles with pagination
func (s *SysRoleAppService) GetSysRolePage(ctx context.Context, query *domain.SysRolePageQuery) (*entity.PaginationResult[*domain.SysRole], error) {
	result, err := s.repo.FindPage(ctx, query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ChangeRoleStatus changes role status
func (s *SysRoleAppService) ChangeRoleStatus(ctx context.Context, cmd ChangeRoleStatusCommand) error {
	role, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	return role.ChangeStatus(cmd.Status, cmd.UserID)
}
