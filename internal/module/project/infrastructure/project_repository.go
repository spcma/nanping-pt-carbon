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

type ProjectRepository struct {
	*db.BaseRepository[domain.Project]
}

func NewProjectRepository(_db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{
		BaseRepository: db.NewBaseRepository[domain.Project](_db),
	}
}

func (r *ProjectRepository) Create(ctx context.Context, project *domain.Project) error {
	return r.GetDB(ctx).WithContext(ctx).Create(project).Error
}

func (r *ProjectRepository) Update(ctx context.Context, project *domain.Project) error {
	return r.GetDB(ctx).WithContext(ctx).Save(project).Error
}

func (r *ProjectRepository) UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.GetDB(ctx).WithContext(ctx).Model(&domain.Project{}).Where("id = ? AND "+entity.FieldDeleteBy+" = 0", id).Updates(updates).Error
}

func (r *ProjectRepository) Delete(ctx context.Context, id, uid int64) error {
	// 实现逻辑删除，更新状态而不是物理删除
	updates := map[string]interface{}{
		"deleteBy":   uid,
		"deleteTime": timeutil.Now(),
	}
	return r.GetDB(ctx).WithContext(ctx).Model(&domain.Project{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ProjectRepository) FindByID(ctx context.Context, id int64) (*domain.Project, error) {
	var project domain.Project
	err := r.GetDB(ctx).WithContext(ctx).
		Table("project").
		Where("id = ?", id).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) FindByCode(ctx context.Context, code string) (*domain.Project, error) {
	var project domain.Project
	err := r.GetDB(ctx).WithContext(ctx).
		Table("project").
		Where("code = ?", code).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) FindList(ctx context.Context) ([]*domain.Project, error) {
	var projects []*domain.Project
	err := r.GetDB(ctx).WithContext(ctx).Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) FindPage(ctx context.Context, query *domain.ProjectPageQuery) ([]*domain.Project, int64, error) {
	// 使用通用分页助手
	helper := db.NewPaginationHelper[*domain.Project](r.GetDB(ctx))
	result, err := helper.PageQuery(query.PageNum, query.PageSize, func(dq *gorm.DB) *gorm.DB {
		// 构建基础查询 - 使用 delete_by 条件
		dq = dq.WithContext(ctx).
			Table("project").
			Where(entity.FieldDeleteBy + " = 0")

		// 动态构建查询条件
		if query.Name != "" {
			dq = dq.Where("name LIKE ?", "%"+query.Name+"%")
		}
		if query.Code != "" {
			dq = dq.Where("code = ?", query.Code)
		}
		if query.Status != "" {
			dq = dq.Where("status = ?", query.Status)
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
		return nil, 0, err
	}
	return result.Data, result.Total, nil
}

func (r *ProjectRepository) FindListByStatus(ctx context.Context, status domain.ProjectStatus) ([]*domain.Project, error) {
	var projects []*domain.Project
	err := r.GetDB(ctx).WithContext(ctx).
		Table("project").
		Where("status = ? AND "+entity.FieldDeleteBy+" = 0", status).
		Find(&projects).Error
	return projects, err
}
