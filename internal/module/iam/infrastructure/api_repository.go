package infrastructure

import (
	"app/internal/module/iam/domain"
	"app/internal/shared/db"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"

	"gorm.io/gorm"
)

type ApiRepository struct {
}

func NewApiRepository(_db *gorm.DB) *ApiRepository {
	return &ApiRepository{}
}

func (u *ApiRepository) Create(ctx context.Context, api *domain.Api) error {
	return db.GetDB(ctx).WithContext(ctx).Create(api).Error
}

func (u *ApiRepository) Update(ctx context.Context, api *domain.Api) error {
	return db.GetDB(ctx).WithContext(ctx).Updates(api).Error
}

func (u *ApiRepository) UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error {
	return db.GetDB(ctx).WithContext(ctx).Model(&domain.Api{}).Where("id = ? AND "+entity.FieldDeleteBy+" = 0", id).Updates(updates).Error
}

func (u *ApiRepository) Delete(ctx context.Context, id int64) error {
	// 实现逻辑删除，更新状态而不是物理删除
	updates := map[string]interface{}{
		"deleteBy":   0,
		"deleteTime": timeutil.Now(),
	}
	return db.GetDB(ctx).WithContext(ctx).Model(&domain.Api{}).Where("id = ?", id).Updates(updates).Error
}

func (u *ApiRepository) FindByID(ctx context.Context, id int64) (*domain.Api, error) {
	var api domain.Api
	err := db.GetDB(ctx).WithContext(ctx).
		Table("sys_api").
		Where("id = ?", id).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&api).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &api, nil
}

func (u *ApiRepository) FindByUsername(ctx context.Context, username string) (*domain.Api, error) {
	var api domain.Api
	err := db.GetDB(ctx).WithContext(ctx).
		Table("sys_api").
		Where("username = ?", username).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&api).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &api, nil
}

func (u *ApiRepository) FindList(ctx context.Context) ([]*domain.Api, error) {
	var apis []*domain.Api
	err := db.GetDB(ctx).WithContext(ctx).Find(&apis).Error
	return apis, err
}

func (u *ApiRepository) FindPage(ctx context.Context, pageNum, pageSize int64, name string) ([]*domain.Api, int64, error) {
	// 使用通用分页助手
	helper := db.NewPaginationHelper[*domain.Api](db.GetDB(ctx))
	result, err := helper.PageQuery(int(pageNum), int(pageSize), func(dq *gorm.DB) *gorm.DB {
		dq = dq.WithContext(ctx).
			Table("sys_api").
			Where(entity.FieldDeleteBy + " = 0")

		// 动态构建查询条件
		if name != "" {
			dq = dq.Where("name LIKE ?", "%"+name+"%")
		}

		return dq
	})
	if err != nil {
		return nil, 0, err
	}
	return result.Data, result.Total, nil
}

func (u *ApiRepository) FindAll(ctx context.Context) ([]*domain.Api, error) {
	var apis []*domain.Api
	err := db.GetDB(ctx).WithContext(ctx).
		Table("sys_api").
		Where(entity.FieldDeleteBy + " = 0").
		Find(&apis).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return apis, err
}

func (u *ApiRepository) FindByCode(ctx context.Context, code string) (*domain.Api, error) {
	var api domain.Api
	err := db.GetDB(ctx).WithContext(ctx).
		Table("sys_api").
		Where("code = ?", code).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&api).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &api, nil
}
