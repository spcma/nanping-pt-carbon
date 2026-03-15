package application

import (
	"app/internal/module/iam/domain"
	"context"
)

// CreateSysApiCommand create system API command
type CreateSysApiCommand struct {
	Name       string `json:"name"`
	Code       string `json:"code"`
	Uri        string `json:"uri"`
	MethodType string `json:"methodType"`
	Status     string `json:"status"`
	UserID     int64  `json:"userId"`
}

// UpdateSysApiCommand update system API command
type UpdateSysApiCommand struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Uri        string `json:"uri"`
	MethodType string `json:"methodType"`
	Status     string `json:"status"`
	UserID     int64  `json:"userId"`
}

// SysApiAppService system API application service
type SysApiAppService struct {
	repo domain.SysApiRepository
}

// NewSysApiAppService creates system API application service
func NewSysApiAppService(repo domain.SysApiRepository) *SysApiAppService {
	return &SysApiAppService{repo: repo}
}

// CreateSysApi creates a system API
func (s *SysApiAppService) CreateSysApi(ctx context.Context, cmd CreateSysApiCommand) (int64, error) {
	api, err := domain.NewSysApi(cmd.Name, cmd.Code, cmd.Uri, cmd.MethodType, cmd.Status, cmd.UserID)
	if err != nil {
		return 0, err
	}
	err = s.repo.Create(ctx, api)
	if err != nil {
		return 0, err
	}
	return api.Id, nil
}

// UpdateSysApi updates a system API
func (s *SysApiAppService) UpdateSysApi(ctx context.Context, cmd UpdateSysApiCommand) error {
	api, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	return api.UpdateInfo(cmd.Name, cmd.Uri, cmd.MethodType, cmd.Status, cmd.UserID)
}

// DeleteSysApi deletes a system API
func (s *SysApiAppService) DeleteSysApi(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// GetSysAPIByID gets system API by ID
func (s *SysApiAppService) GetSysAPIByID(ctx context.Context, id int64) (*domain.SysApi, error) {
	return s.repo.FindByID(ctx, id)
}

// GetSysAPIByCode gets system API by code
func (s *SysApiAppService) GetSysAPIByCode(ctx context.Context, code string) (*domain.SysApi, error) {
	return s.repo.FindByCode(ctx, code)
}

// GetSysAPIPage queries system APIs with pagination
func (s *SysApiAppService) GetSysAPIPage(ctx context.Context, pageNum, pageSize int64, name string) ([]*domain.SysApi, int64, error) {
	return s.repo.FindPage(ctx, pageNum, pageSize, name)
}

// GetAllSysAPIs gets all system APIs
func (s *SysApiAppService) GetAllSysAPIs(ctx context.Context) ([]*domain.SysApi, error) {
	return s.repo.FindAll(ctx)
}
