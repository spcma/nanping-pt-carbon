package persistence

import (
	"app/internal/module/project/domain"
	"app/internal/shared/db"
	"app/internal/shared/timeutil"
	"context"
	"strings"

	"gorm.io/gorm"
)

type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates project repository
func NewProjectRepository(db *gorm.DB) domain.ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(ctx context.Context, project *domain.Project) error {
	return db.WithContext(ctx, r.db).Create(project).Error
}

func (r *projectRepository) Update(ctx context.Context, project *domain.Project) error {
	return db.WithContext(ctx, r.db).Save(project).Error
}

func (r *projectRepository) Delete(ctx context.Context, id int64, userID int64) error {
	// 逻辑删除：只更新 delete_user 和 delete_time 字段
	updates := map[string]interface{}{
		"delete_by":   userID,
		"delete_time": timeutil.New(),
	}
	return db.WithContext(ctx, r.db).Model(&domain.Project{}).Where("id = ?", id).Updates(updates).Error
}

func (r *projectRepository) FindByID(ctx context.Context, id int64) (*domain.Project, error) {
	var project domain.Project
	err := db.WithContext(ctx, r.db).Where("id = ? AND delete_time IS NULL", id).First(&project).Error
	return &project, err
}

func (r *projectRepository) FindByCode(ctx context.Context, code string) (*domain.Project, error) {
	var project domain.Project
	err := db.WithContext(ctx, r.db).Where("code = ? AND delete_time IS NULL", code).First(&project).Error
	return &project, err
}

func (r *projectRepository) FindPage(ctx context.Context, query domain.ProjectPageQuery) ([]*domain.Project, int64, error) {
	var projects []*domain.Project
	var total int64

	tx := db.WithContext(ctx, r.db).Model(&domain.Project{}).Where("delete_time IS NULL")

	// 动态构建查询条件
	if query.Name != "" {
		tx = tx.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Code != "" {
		tx = tx.Where("code = ?", query.Code)
	}
	if query.Status != "" {
		tx = tx.Where("status = ?", query.Status)
	}

	// 计数
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	if query.SortBy != "" {
		order := query.SortBy + " " + strings.ToUpper(query.SortOrder)
		tx = tx.Order(order)
	} else {
		tx = tx.Order("id DESC") // 默认按 ID 降序
	}

	// 分页
	offset := (query.PageNum - 1) * query.PageSize
	err := tx.Offset(int(offset)).Limit(int(query.PageSize)).Find(&projects).Error
	return projects, total, err
}

func (r *projectRepository) FindListByStatus(ctx context.Context, status domain.ProjectStatus) ([]*domain.Project, error) {
	var projects []*domain.Project
	err := db.WithContext(ctx, r.db).Where("status = ? AND delete_time IS NULL", status).Find(&projects).Error
	return projects, err
}
