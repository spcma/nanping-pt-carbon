package infrastructure

import (
	"app/internal/module/methodology/domain"
	"app/internal/shared/db"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"

	"gorm.io/gorm"
)

type MethodologyRepository struct {
	*db.BaseRepository[domain.Methodology]
}

func NewMethodologyRepository(_db *gorm.DB) *MethodologyRepository {
	return &MethodologyRepository{
		BaseRepository: db.NewBaseRepository[domain.Methodology](_db),
	}
}

func (r *MethodologyRepository) Create(ctx context.Context, methodology *domain.Methodology) error {
	return r.GetDB(ctx).WithContext(ctx).Create(methodology).Error
}

func (r *MethodologyRepository) Update(ctx context.Context, methodology *domain.Methodology) error {
	return r.GetDB(ctx).WithContext(ctx).Save(methodology).Error
}

func (r *MethodologyRepository) UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.GetDB(ctx).WithContext(ctx).Model(&domain.Methodology{}).Where("id = ? AND "+entity.FieldDeleteBy+" = 0", id).Updates(updates).Error
}

func (r *MethodologyRepository) Delete(ctx context.Context, id, uid int64) error {
	// 实现逻辑删除，更新状态而不是物理删除
	updates := map[string]interface{}{
		"deleteBy":   uid,
		"deleteTime": timeutil.Now(),
	}
	return r.GetDB(ctx).WithContext(ctx).Model(&domain.Methodology{}).Where("id = ?", id).Updates(updates).Error
}

func (r *MethodologyRepository) FindByID(ctx context.Context, id int64) (*domain.Methodology, error) {
	var methodology domain.Methodology
	err := r.GetDB(ctx).WithContext(ctx).
		Table("methodology").
		Where("id = ?", id).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&methodology).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &methodology, nil
}

func (r *MethodologyRepository) FindByQuery(ctx context.Context, query *domain.MethodologyQuery) (*domain.Methodology, error) {
	var methodology domain.Methodology
	err := r.GetDB(ctx).WithContext(ctx).
		Table("methodology").
		Where("code = ?", query.Code).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&methodology).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &methodology, nil
}

func (r *MethodologyRepository) FindList(ctx context.Context) ([]*domain.Methodology, error) {
	var methodologies []*domain.Methodology
	err := r.GetDB(ctx).WithContext(ctx).Find(&methodologies).Error
	return methodologies, err
}

func (r *MethodologyRepository) FindPage(ctx context.Context, query *domain.MethodologyPageQuery) (*entity.PaginationResult[*domain.Methodology], error) {
	// 使用通用分页助手
	helper := db.NewPaginationHelper[*domain.Methodology](r.GetDB(ctx))
	result, err := helper.PageQuery(int(query.PageNum), int(query.PageSize), func(dq *gorm.DB) *gorm.DB {
		// 构建基础查询 - 使用 delete_by 条件
		dq = dq.WithContext(ctx).
			Table("methodology").
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
	return result, nil
}

func (r *MethodologyRepository) FindListByStatus(ctx context.Context, status domain.MethodologyStatus) ([]*domain.Methodology, error) {
	var methodologies []*domain.Methodology
	err := r.GetDB(ctx).WithContext(ctx).
		Table("methodology").
		Where("status = ? AND "+entity.FieldDeleteBy+" = 0", status).
		Find(&methodologies).Error
	return methodologies, err
}
