package project

import (
	"app/internal/shared/db"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"

	"gorm.io/gorm"
)

type ProjectRepository struct {
}

func NewProjectRepository(_db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{}
}

func (r *ProjectRepository) Create(ctx context.Context, project *Project) error {
	return db.GetDB(ctx).WithContext(ctx).Create(project).Error
}

func (r *ProjectRepository) Update(ctx context.Context, project *Project) error {
	return db.GetDB(ctx).WithContext(ctx).Save(project).Error
}

func (r *ProjectRepository) UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error {
	return db.GetDB(ctx).WithContext(ctx).Model(&Project{}).Where("id = ? AND "+entity.FieldDeleteBy+" = 0", id).Updates(updates).Error
}

func (r *ProjectRepository) Delete(ctx context.Context, id, uid int64) error {
	// 实现逻辑删除，更新状态而不是物理删除
	updates := map[string]interface{}{
		"deleteBy":   uid,
		"deleteTime": timeutil.Now(),
	}
	return db.GetDB(ctx).WithContext(ctx).Model(&Project{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ProjectRepository) FindByID(ctx context.Context, id int64) (*Project, error) {
	var project Project
	err := db.GetDB(ctx).WithContext(ctx).
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

func (r *ProjectRepository) FindByCode(ctx context.Context, code string) (*Project, error) {
	var project Project
	err := db.GetDB(ctx).WithContext(ctx).
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

// FindByQuery 根据条件查询项目（支持多条件组合）
func (r *ProjectRepository) FindByQuery(ctx context.Context, query *ProjectQuery) (*Project, error) {
	txDB := db.GetDB(ctx).WithContext(ctx).
		Table("project").
		Where(entity.FieldDeleteBy + " = 0")

	// 动态构建 WHERE 条件
	if query.ID > 0 {
		txDB = txDB.Where("id = ?", query.ID)
	}
	if query.Code != "" {
		txDB = txDB.Where("code = ?", query.Code)
	}
	if query.Name != "" {
		txDB = txDB.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Status != "" {
		txDB = txDB.Where("status = ?", query.Status)
	}

	var project Project
	err := txDB.Take(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) FindList(ctx context.Context) ([]*Project, error) {
	var projects []*Project
	err := db.GetDB(ctx).WithContext(ctx).Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) FindPage(ctx context.Context, query *ProjectPageQuery) (*entity.PaginationResult[*Project], error) {
	helper := db.NewPaginationHelper[*Project](db.GetDB(ctx))

	pageQuery, err := helper.PageQuery(query.PageNum, query.PageSize, func(dq *gorm.DB) *gorm.DB {
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
		return nil, err
	}
	return pageQuery, nil
}

func (r *ProjectRepository) FindListByStatus(ctx context.Context, status ProjectStatus) ([]*Project, error) {
	var projects []*Project
	err := db.GetDB(ctx).WithContext(ctx).
		Table("project").
		Where("status = ? AND "+entity.FieldDeleteBy+" = 0", status).
		Find(&projects).Error
	return projects, err
}
