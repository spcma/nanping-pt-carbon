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
	*db.BaseRepository[domain.SysApi]
}

func NewApiRepository(_db *gorm.DB) *ApiRepository {
	return &ApiRepository{
		BaseRepository: db.NewBaseRepository[domain.SysApi](_db),
	}
}

func (u *ApiRepository) Create(ctx context.Context, user *domain.SysApi) error {
	return u.GetDB(ctx).WithContext(ctx).Create(user).Error
}

func (u *ApiRepository) Update(ctx context.Context, user *domain.SysApi) error {
	return u.GetDB(ctx).WithContext(ctx).Updates(user).Error
}

func (u *ApiRepository) UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error {
	return u.GetDB(ctx).WithContext(ctx).Model(&domain.SysApi{}).Where("id = ? AND "+entity.FieldDeleteBy+" = 0", id).Updates(updates).Error
}

func (u *ApiRepository) Delete(ctx context.Context, id int64) error {
	// 实现逻辑删除，更新状态而不是物理删除
	updates := map[string]interface{}{
		"deleteBy":   0,
		"deleteTime": timeutil.Now(),
	}
	return u.GetDB(ctx).WithContext(ctx).Model(&domain.SysApi{}).Where("id = ?", id).Updates(updates).Error
}

func (u *ApiRepository) FindByID(ctx context.Context, id int64) (*domain.SysApi, error) {
	var user domain.SysApi
	err := u.GetDB(ctx).WithContext(ctx).
		Table("apis").
		Where("id = ?", id).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (u *ApiRepository) FindByUsername(ctx context.Context, username string) (*domain.SysApi, error) {
	var user domain.SysApi
	err := u.GetDB(ctx).WithContext(ctx).
		Table("apis").
		Where("username = ?", username).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (u *ApiRepository) FindList(ctx context.Context) ([]*domain.SysApi, error) {
	var users []*domain.SysApi
	err := u.GetDB(ctx).WithContext(ctx).Find(&users).Error
	return users, err
}

func (u *ApiRepository) FindPage(ctx context.Context, pageNum, pageSize int64, name string) ([]*domain.SysApi, int64, error) {
	// 使用通用分页助手
	helper := db.NewPaginationHelper[*domain.SysApi](u.GetDB(ctx))
	result, err := helper.PageQuery(int(pageNum), int(pageSize), func(dq *gorm.DB) *gorm.DB {
		// 构建基础查询 - 使用 delete_by 条件
		dq = dq.WithContext(ctx).
			Table("apis").
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

func (u *ApiRepository) FindAll(ctx context.Context) ([]*domain.SysApi, error) {
	var apis []*domain.SysApi
	err := u.GetDB(ctx).WithContext(ctx).
		Table("apis").
		Where(entity.FieldDeleteBy + " = 0").
		Find(&apis).Error
	return apis, err
}

func (u *ApiRepository) FindByCode(ctx context.Context, code string) (*domain.SysApi, error) {
	var api domain.SysApi
	err := u.GetDB(ctx).WithContext(ctx).
		Table("apis").
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
