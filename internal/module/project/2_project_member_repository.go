package project

import (
	"app/internal/shared/db"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"

	"gorm.io/gorm"
)

type ProjectMembersRepository struct{}

func NewProjectMembersRepository(_db *gorm.DB) ProjectMembersRepo {
	return &ProjectMembersRepository{}
}

func (r *ProjectMembersRepository) Create(ctx context.Context, projectMembers *ProjectMembers) error {
	return db.GetDB(ctx).WithContext(ctx).Create(projectMembers).Error
}

func (r *ProjectMembersRepository) Update(ctx context.Context, projectMembers *ProjectMembers) error {
	return db.GetDB(ctx).WithContext(ctx).Save(projectMembers).Error
}

func (r *ProjectMembersRepository) Delete(ctx context.Context, id, uid int64) error {
	// 实现逻辑删除，更新状态而不是物理删除
	updates := map[string]interface{}{
		"deleteBy":   uid,
		"deleteTime": timeutil.Now(),
	}
	return db.GetDB(ctx).WithContext(ctx).Model(&ProjectMembers{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ProjectMembersRepository) FindByID(ctx context.Context, id int64) (*ProjectMembers, error) {
	var projectMember ProjectMembers
	err := db.GetDB(ctx).WithContext(ctx).
		Table("project_members").
		Where("id = ?", id).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&projectMember).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &projectMember, nil
}

func (r *ProjectMembersRepository) FindByProjectID(ctx context.Context, projectID int64) ([]*ProjectMembers, error) {
	var projectMembers []*ProjectMembers
	err := db.GetDB(ctx).WithContext(ctx).
		Table("project_members").
		Where("project_id = ? AND "+entity.FieldDeleteBy+" = 0", projectID).
		Find(&projectMembers).Error
	return projectMembers, err
}

func (r *ProjectMembersRepository) FindByUserID(ctx context.Context, userID int64) ([]*ProjectMembers, error) {
	var projectMembers []*ProjectMembers
	err := db.GetDB(ctx).WithContext(ctx).
		Table("project_members").
		Where("user_id = ? AND "+entity.FieldDeleteBy+" = 0", userID).
		Find(&projectMembers).Error
	return projectMembers, err
}

func (r *ProjectMembersRepository) FindByProjectAndUser(ctx context.Context, projectID, userID int64) (*ProjectMembers, error) {
	var projectMember ProjectMembers
	err := db.GetDB(ctx).WithContext(ctx).
		Table("project_members").
		Where("project_id = ? AND user_id = ? AND "+entity.FieldDeleteBy+" = 0", projectID, userID).
		Take(&projectMember).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &projectMember, nil
}

func (r *ProjectMembersRepository) FindPage(ctx context.Context, query *ProjectMembersPageQuery) (*entity.PaginationResult[*ProjectMembers], error) {
	// 使用通用分页助手
	helper := db.NewPaginationHelper[*ProjectMembers](db.GetDB(ctx))
	result, err := helper.PageQuery(query.PageNum, query.PageSize, func(dq *gorm.DB) *gorm.DB {
		// 构建基础查询 - 使用 delete_by 条件
		dq = dq.WithContext(ctx).
			Table("project_members").
			Where(entity.FieldDeleteBy + " = 0")

		// 动态构建查询条件
		if query.ProjectId != 0 {
			dq = dq.Where("project_id = ?", query.ProjectId)
		}
		if query.UserId != 0 {
			dq = dq.Where("user_id = ?", query.UserId)
		}
		if query.Role != "" {
			dq = dq.Where("role = ?", query.Role)
		}

		// 排序
		if query.SortBy != "" {
			order := query.SortBy + " " + query.SortOrder
			dq = dq.Order(order)
		} else {
			dq = dq.Order("id DESC") // 默认按 ID 降序
		}

		return dq
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
