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

// MethodologyRepositoryImpl 方法学仓储实现
type MethodologyRepositoryImpl struct{}

// NewMethodologyRepository 创建方法学仓储实例
func NewMethodologyRepository() *MethodologyRepositoryImpl {
	return &MethodologyRepositoryImpl{}
}

func (r *MethodologyRepositoryImpl) Create(ctx context.Context, methodology *domain.Methodology) error {
	return db.GetDB(ctx).WithContext(ctx).Create(methodology).Error
}

func (r *MethodologyRepositoryImpl) Update(ctx context.Context, methodology *domain.Methodology) error {
	return db.GetDB(ctx).WithContext(ctx).Save(methodology).Error
}

func (r *MethodologyRepositoryImpl) Delete(ctx context.Context, id, uid int64) error {
	updates := map[string]interface{}{
		"deleteBy":   uid,
		"deleteTime": timeutil.Now(),
	}
	return db.GetDB(ctx).WithContext(ctx).Model(&domain.Methodology{}).Where("id = ?", id).Updates(updates).Error
}

func (r *MethodologyRepositoryImpl) FindByID(ctx context.Context, id int64) (*domain.Methodology, error) {
	var methodology domain.Methodology
	err := db.GetDB(ctx).WithContext(ctx).
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

func (r *MethodologyRepositoryImpl) FindByQuery(ctx context.Context, query *domain.MethodologyQuery) (*domain.Methodology, error) {
	txDB := db.GetDB(ctx).WithContext(ctx).
		Table("methodology").
		Where(entity.FieldDeleteBy + " = 0")

	if query.ID != nil && *query.ID > 0 {
		txDB = txDB.Where("id = ?", *query.ID)
	}
	if query.Code != nil && *query.Code != "" {
		txDB = txDB.Where("code = ?", *query.Code)
	}
	if query.Name != nil && *query.Name != "" {
		txDB = txDB.Where("name LIKE ?", "%"+*query.Name+"%")
	}
	if query.Status != nil {
		txDB = txDB.Where("status = ?", *query.Status)
	}

	var methodology domain.Methodology
	err := txDB.Take(&methodology).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &methodology, nil
}

func (r *MethodologyRepositoryImpl) FindList(ctx context.Context) ([]*domain.Methodology, error) {
	var methodologies []*domain.Methodology
	err := db.GetDB(ctx).WithContext(ctx).Find(&methodologies).Error
	return methodologies, err
}

func (r *MethodologyRepositoryImpl) FindPage(ctx context.Context, query *domain.MethodologyPageQuery) (*entity.PaginationResult[*domain.Methodology], error) {
	helper := db.NewPaginationHelper[*domain.Methodology](db.GetDB(ctx))

	result, err := helper.PageQuery(query.PageNum, query.PageSize, func(dq *gorm.DB) *gorm.DB {
		dq = dq.WithContext(ctx).
			Table("methodology").
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
	return result, nil
}

func (r *MethodologyRepositoryImpl) FindListByStatus(ctx context.Context, status domain.MethodologyStatus) ([]*domain.Methodology, error) {
	var methodologies []*domain.Methodology
	err := db.GetDB(ctx).WithContext(ctx).
		Table("methodology").
		Where("status = ? AND "+entity.FieldDeleteBy+" = 0", status).
		Find(&methodologies).Error
	return methodologies, err
}
