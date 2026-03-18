package application

import (
	"app/internal/module/project/domain"
	"app/internal/shared/entity"
	"context"
)

// ProjectMembersService 项目成员应用服务
type ProjectMembersService struct {
	repo ProjectMembersRepo
}

// NewProjectMembersService 创建项目成员应用服务
func NewProjectMembersService(repo ProjectMembersRepo) *ProjectMembersService {
	return &ProjectMembersService{repo: repo}
}

// CreateProjectMemberParam 创建项目成员命令
type CreateProjectMemberParam struct {
	ProjectId int64  `json:"projectId"`
	UserId    int64  `json:"userId"`
	Role      string `json:"role"` // owner admin member
	CreateBy  int64  // 操作人 ID
}

// CreateProjectMember 创建项目成员
func (s *ProjectMembersService) CreateProjectMember(ctx context.Context, param CreateProjectMemberParam) (int64, error) {
	projectMembers, err := s.repo.FindByUserID(ctx, param.UserId)
	if err != nil {
		return 0, err
	}

	for _, member := range projectMembers {
		if member.ProjectId == param.ProjectId {
			return member.Id, domain.ErrProjectMemberAlreadyExists
		}
	}

	member, err := domain.NewProjectMembers(param.ProjectId, param.UserId, param.CreateBy, param.Role)
	if err != nil {
		return 0, err
	}

	err = s.repo.Create(ctx, member)
	if err != nil {
		return 0, err
	}

	return member.Id, nil
}

// UpdateProjectMemberParam 更新项目成员命令
type UpdateProjectMemberParam struct {
	ID       int64  `json:"id"`
	Role     string `json:"role"`
	CreateBy int64  // 操作人 ID
}

// UpdateProjectMember 更新项目成员
func (s *ProjectMembersService) UpdateProjectMember(ctx context.Context, param UpdateProjectMemberParam) error {
	member, err := s.repo.FindByID(ctx, param.ID)
	if err != nil {
		return err
	}
	if member == nil {
		return domain.ErrProjectMemberNotFound
	}
	return member.UpdateRole(param.Role, param.CreateBy)
}

// DeleteProjectMemberParam 删除项目成员命令
type DeleteProjectMemberParam struct {
	ID       int64 `json:"id"`
	CreateBy int64
}

// DeleteProjectMember 删除项目成员
func (s *ProjectMembersService) DeleteProjectMember(ctx context.Context, param DeleteProjectMemberParam) error {
	member, err := s.repo.FindByID(ctx, param.ID)
	if err != nil {
		return err
	}
	if member == nil {
		return domain.ErrProjectMemberNotFound
	}
	return member.Delete(param.CreateBy)
}

// GetProjectMemberByID 根据 ID 获取项目成员
func (s *ProjectMembersService) GetProjectMemberByID(ctx context.Context, id int64) (*domain.ProjectMembers, error) {
	return s.repo.FindByID(ctx, id)
}

// GetProjectMembersByProjectID 根据项目 ID 获取成员列表
func (s *ProjectMembersService) GetProjectMembersByProjectID(ctx context.Context, projectID int64) ([]*domain.ProjectMembers, error) {
	return s.repo.FindByProjectID(ctx, projectID)
}

// GetProjectMembersByUserID 根据用户 ID 获取参与的项目
func (s *ProjectMembersService) GetProjectMembersByUserID(ctx context.Context, userID int64) ([]*domain.ProjectMembers, error) {
	return s.repo.FindByUserID(ctx, userID)
}

// GetProjectMemberByProjectAndUser 根据项目和用户获取成员信息
func (s *ProjectMembersService) GetProjectMemberByProjectAndUser(ctx context.Context, projectID, userID int64) (*domain.ProjectMembers, error) {
	return s.repo.FindByProjectAndUser(ctx, projectID, userID)
}

// GetProjectMemberPage 分页查询项目成员
func (s *ProjectMembersService) GetProjectMemberPage(ctx context.Context, query *domain.ProjectMembersPageQuery) (*entity.PaginationResult[*domain.ProjectMembers], error) {
	return s.repo.FindPage(ctx, query)
}
