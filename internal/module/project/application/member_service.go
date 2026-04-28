package application

import (
	"app/internal/module/project/domain"
	"app/internal/shared/entity"
	"context"
	"errors"
)

// ProjectMembersAppService 项目成员应用服务
type ProjectMembersAppService struct {
	repo domain.ProjectMembersRepository
}

// NewProjectMembersAppService 创建项目成员应用服务
func NewProjectMembersAppService(repo domain.ProjectMembersRepository) *ProjectMembersAppService {
	return &ProjectMembersAppService{repo: repo}
}

// CreateProjectMemberCommand 创建项目成员命令
type CreateProjectMemberCommand struct {
	ProjectId int64  `json:"projectId"`
	UserId    int64  `json:"userId"`
	Role      string `json:"role"` // owner admin member
	CreateBy  int64  // 操作人 ID
}

// CreateProjectMember 创建项目成员
func (s *ProjectMembersAppService) CreateProjectMember(ctx context.Context, cmd *CreateProjectMemberCommand) (int64, error) {
	// 检查成员是否已存在
	existing, err := s.repo.FindByProjectAndUser(ctx, cmd.ProjectId, cmd.UserId)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return 0, errors.New("该用户已是项目成员")
	}

	// 调用领域层创建聚合根
	member, err := domain.NewProjectMembers(cmd.ProjectId, cmd.UserId, cmd.CreateBy, cmd.Role)
	if err != nil {
		return 0, err
	}

	// 持久化
	err = s.repo.Create(ctx, member)
	if err != nil {
		return 0, err
	}

	return member.Id, nil
}

// UpdateProjectMemberCommand 更新项目成员命令
type UpdateProjectMemberCommand struct {
	ID       int64  `json:"id"`
	Role     string `json:"role"`
	CreateBy int64  // 操作人 ID
}

// UpdateProjectMember 更新项目成员
func (s *ProjectMembersAppService) UpdateProjectMember(ctx context.Context, cmd UpdateProjectMemberCommand) error {
	// 获取聚合根
	member, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if member == nil {
		return errors.New("项目成员不存在")
	}

	// 调用领域行为
	return member.UpdateRole(cmd.Role, cmd.CreateBy)
}

// DeleteProjectMemberCommand 删除项目成员命令
type DeleteProjectMemberCommand struct {
	ID       int64 `json:"id"`
	CreateBy int64
}

// DeleteProjectMember 删除项目成员
func (s *ProjectMembersAppService) DeleteProjectMember(ctx context.Context, cmd DeleteProjectMemberCommand) error {
	// 获取聚合根
	member, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if member == nil {
		return errors.New("项目成员不存在")
	}

	// 调用领域行为
	return member.Delete(cmd.CreateBy)
}

// GetProjectMemberByID 根据 ID 获取项目成员
func (s *ProjectMembersAppService) GetProjectMemberByID(ctx context.Context, id int64) (*domain.ProjectMembers, error) {
	return s.repo.FindByID(ctx, id)
}

// GetProjectMembersByProjectID 根据项目 ID 获取成员列表
func (s *ProjectMembersAppService) GetProjectMembersByProjectID(ctx context.Context, projectID int64) ([]*domain.ProjectMembers, error) {
	return s.repo.FindByProjectID(ctx, projectID)
}

// GetProjectMembersByUserID 根据用户 ID 获取参与的项目
func (s *ProjectMembersAppService) GetProjectMembersByUserID(ctx context.Context, userID int64) ([]*domain.ProjectMembers, error) {
	return s.repo.FindByUserID(ctx, userID)
}

// GetProjectMemberByProjectAndUser 根据项目和用户获取成员信息
func (s *ProjectMembersAppService) GetProjectMemberByProjectAndUser(ctx context.Context, projectID, userID int64) (*domain.ProjectMembers, error) {
	return s.repo.FindByProjectAndUser(ctx, projectID, userID)
}

// GetProjectMemberPage 分页查询项目成员
func (s *ProjectMembersAppService) GetProjectMemberPage(ctx context.Context, query *domain.ProjectMembersPageQuery) (*entity.PaginationResult[*domain.ProjectMembers], error) {
	return s.repo.FindPage(ctx, query)
}
