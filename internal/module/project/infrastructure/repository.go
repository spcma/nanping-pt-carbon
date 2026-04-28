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

// ProjectRepositoryImpl 项目仓储实现
type ProjectRepositoryImpl struct {
}

// NewProjectRepository 创建项目仓储实例
func NewProjectRepository() *ProjectRepositoryImpl {
	return &ProjectRepositoryImpl{}
}

func (r *ProjectRepositoryImpl) Create(ctx context.Context, project *domain.Project) error {
	return db.GetDB(ctx).WithContext(ctx).Create(project).Error
}

func (r *ProjectRepositoryImpl) Update(ctx context.Context, project *domain.Project) error {
	return db.GetDB(ctx).WithContext(ctx).Save(project).Error
}

func (r *ProjectRepositoryImpl) Delete(ctx context.Context, id, uid int64) error {
	updates := map[string]interface{}{
		"deleteBy":   uid,
		"deleteTime": timeutil.Now(),
	}
	return db.GetDB(ctx).WithContext(ctx).Model(&domain.Project{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ProjectRepositoryImpl) FindByID(ctx context.Context, id int64) (*domain.Project, error) {
	var project domain.Project
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

func (r *ProjectRepositoryImpl) FindByCode(ctx context.Context, code string) (*domain.Project, error) {
	var project domain.Project
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

func (r *ProjectRepositoryImpl) FindByQuery(ctx context.Context, query *domain.ProjectQuery) (*domain.Project, error) {
	txDB := db.GetDB(ctx).WithContext(ctx).
		Table("project").
		Where(entity.FieldDeleteBy + " = 0")

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

	var project domain.Project
	err := txDB.Take(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepositoryImpl) FindList(ctx context.Context) ([]*domain.Project, error) {
	var projects []*domain.Project
	err := db.GetDB(ctx).WithContext(ctx).Find(&projects).Error
	return projects, err
}

func (r *ProjectRepositoryImpl) FindPage(ctx context.Context, query *domain.ProjectPageQuery) (*entity.PaginationResult[*domain.Project], error) {
	helper := db.NewPaginationHelper[*domain.Project](db.GetDB(ctx))

	pageQuery, err := helper.PageQuery(query.PageNum, query.PageSize, func(dq *gorm.DB) *gorm.DB {
		dq = dq.WithContext(ctx).
			Table("project").
			Where(entity.FieldDeleteBy + " = 0")

		if query.Name != "" {
			dq = dq.Where("name LIKE ?", "%"+query.Name+"%")
		}
		if query.Code != "" {
			dq = dq.Where("code = ?", query.Code)
		}
		if query.Status != "" {
			dq = dq.Where("status = ?", query.Status)
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
	return pageQuery, nil
}

func (r *ProjectRepositoryImpl) FindListByStatus(ctx context.Context, status domain.ProjectStatus) ([]*domain.Project, error) {
	var projects []*domain.Project
	err := db.GetDB(ctx).WithContext(ctx).
		Table("project").
		Where("status = ? AND "+entity.FieldDeleteBy+" = 0", status).
		Find(&projects).Error
	return projects, err
}
