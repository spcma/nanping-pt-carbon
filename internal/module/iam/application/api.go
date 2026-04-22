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

// ApiAppService system API application service
type ApiAppService struct {
	repo ApiRepository
}

// NewApiAppService creates system API application service
func NewApiAppService(repo ApiRepository) *ApiAppService {
	return &ApiAppService{repo: repo}
}

// CreateApi creates a system API
func (s *ApiAppService) CreateApi(ctx context.Context, cmd CreateSysApiCommand) (int64, error) {
	api, err := domain.NewApi(cmd.Name, cmd.Code, cmd.Uri, cmd.MethodType, cmd.Status, cmd.UserID)
	if err != nil {
		return 0, err
	}
	err = s.repo.Create(ctx, api)
	if err != nil {
		return 0, err
	}
	return api.Id, nil
}

// UpdateApi updates a system API
func (s *ApiAppService) UpdateApi(ctx context.Context, cmd UpdateSysApiCommand) error {
	api, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	return api.UpdateInfo(cmd.Name, cmd.Uri, cmd.MethodType, cmd.Status, cmd.UserID)
}

// DeleteApi deletes a system API
func (s *ApiAppService) DeleteApi(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// GetAPIByID gets system API by ID
func (s *ApiAppService) GetAPIByID(ctx context.Context, id int64) (*domain.Api, error) {
	return s.repo.FindByID(ctx, id)
}

// GetAPIByCode gets system API by code
func (s *ApiAppService) GetAPIByCode(ctx context.Context, code string) (*domain.Api, error) {
	return s.repo.FindByCode(ctx, code)
}

// GetAPIPage queries system APIs with pagination
func (s *ApiAppService) GetAPIPage(ctx context.Context, pageNum, pageSize int64, name string) ([]*domain.Api, int64, error) {
	return s.repo.FindPage(ctx, pageNum, pageSize, name)
}

// GetAllAPIs gets all system APIs
func (s *ApiAppService) GetAllAPIs(ctx context.Context) ([]*domain.Api, error) {
	return s.repo.FindAll(ctx)
}
