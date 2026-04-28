package infrastructure

import (
	"app/internal/module/project/domain"
	"app/internal/shared/db"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"

	"gorm.io/gorm"
)

// ProjectMembersRepositoryImpl 项目成员仓储实现
type ProjectMembersRepositoryImpl struct{}

// NewProjectMembersRepository 创建项目成员仓储实例
func NewProjectMembersRepository() *ProjectMembersRepositoryImpl {
	return &ProjectMembersRepositoryImpl{}
}

func (r *ProjectMembersRepositoryImpl) Create(ctx context.Context, member *domain.ProjectMembers) error {
	return db.GetDB(ctx).WithContext(ctx).Create(member).Error
}

func (r *ProjectMembersRepositoryImpl) Update(ctx context.Context, member *domain.ProjectMembers) error {
	return db.GetDB(ctx).WithContext(ctx).Save(member).Error
}

func (r *ProjectMembersRepositoryImpl) Delete(ctx context.Context, id, uid int64) error {
	updates := map[string]interface{}{
		"deleteBy":   uid,
		"deleteTime": timeutil.Now(),
	}
	return db.GetDB(ctx).WithContext(ctx).Model(&domain.ProjectMembers{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ProjectMembersRepositoryImpl) FindByID(ctx context.Context, id int64) (*domain.ProjectMembers, error) {
	var member domain.ProjectMembers
	err := db.GetDB(ctx).WithContext(ctx).
		Table("project_members").
		Where("id = ?", id).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

func (r *ProjectMembersRepositoryImpl) FindByProjectID(ctx context.Context, projectID int64) ([]*domain.ProjectMembers, error) {
	var members []*domain.ProjectMembers
	err := db.GetDB(ctx).WithContext(ctx).
		Table("project_members").
		Where("project_id = ? AND "+entity.FieldDeleteBy+" = 0", projectID).
		Find(&members).Error
	return members, err
}

func (r *ProjectMembersRepositoryImpl) FindByUserID(ctx context.Context, userID int64) ([]*domain.ProjectMembers, error) {
	var members []*domain.ProjectMembers
	err := db.GetDB(ctx).WithContext(ctx).
		Table("project_members").
		Where("user_id = ? AND "+entity.FieldDeleteBy+" = 0", userID).
		Find(&members).Error
	return members, err
}

func (r *ProjectMembersRepositoryImpl) FindByProjectAndUser(ctx context.Context, projectID, userID int64) (*domain.ProjectMembers, error) {
	var member domain.ProjectMembers
	err := db.GetDB(ctx).WithContext(ctx).
		Table("project_members").
		Where("project_id = ? AND user_id = ? AND "+entity.FieldDeleteBy+" = 0", projectID, userID).
		Take(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

func (r *ProjectMembersRepositoryImpl) FindPage(ctx context.Context, query *domain.ProjectMembersPageQuery) (*entity.PaginationResult[*domain.ProjectMembers], error) {
	helper := db.NewPaginationHelper[*domain.ProjectMembers](db.GetDB(ctx))

	result, err := helper.PageQuery(query.PageNum, query.PageSize, func(dq *gorm.DB) *gorm.DB {
		dq = dq.WithContext(ctx).
			Table("project_members").
			Where(entity.FieldDeleteBy + " = 0")

		if query.ProjectId != 0 {
			dq = dq.Where("project_id = ?", query.ProjectId)
		}
		if query.UserId != 0 {
			dq = dq.Where("user_id = ?", query.UserId)
		}
		if query.Role != "" {
			dq = dq.Where("role = ?", query.Role)
		}

		if query.SortBy != "" {
			order := query.SortBy + " " + query.SortOrder
			dq = dq.Order(order)
		} else {
			dq = dq.Order("id DESC")
		}

		return dq
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
