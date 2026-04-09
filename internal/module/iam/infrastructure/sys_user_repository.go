package infrastructure

import (
	"app/internal/module/iam/domain"
	"app/internal/shared/db"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"

	"gorm.io/gorm"
)

type UserRepository struct {
	*db.BaseRepository[domain.Users]
}

func NewUserRepository(_db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: db.NewBaseRepository[domain.Users](_db),
	}
}

func (u *UserRepository) Create(ctx context.Context, user *domain.Users) error {
	return u.GetDB(ctx).WithContext(ctx).Create(user).Error
}

func (u *UserRepository) Update(ctx context.Context, user *domain.Users) error {
	return u.GetDB(ctx).WithContext(ctx).Updates(user).Error
}

func (u *UserRepository) UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error {
	return u.GetDB(ctx).WithContext(ctx).Model(&domain.Users{}).Where("id = ? AND "+entity.FieldDeleteBy+" = 0", id).Updates(updates).Error
}

func (u *UserRepository) Delete(ctx context.Context, id, uid int64) error {
	// 实现逻辑删除，更新状态而不是物理删除
	updates := map[string]interface{}{
		"deleteBy":   uid,
		"deleteTime": timeutil.Now(),
	}
	return u.GetDB(ctx).WithContext(ctx).Model(&domain.Users{}).Where("id = ?", id).Updates(updates).Error
}

func (u *UserRepository) FindByID(ctx context.Context, id int64) (*domain.Users, error) {
	var user domain.Users
	err := u.GetDB(ctx).WithContext(ctx).
		Table("users").
		Where("id = ?", id).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&user).Error
	if err != nil {
		return nil, db.ErrFilter(err)
	}
	return &user, nil
}

func (u *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.Users, error) {
	var user domain.Users
	err := u.GetDB(ctx).WithContext(ctx).
		Table("users").
		Where("username = ?", username).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&user).Error

	if err != nil {
		return nil, db.ErrFilter(err)
	}

	return &user, nil
}

func (u *UserRepository) FindByQuery(ctx context.Context, query *domain.UserQuery) (*domain.Users, error) {
	tx := u.GetDB(ctx).WithContext(ctx).Table("users").Where(entity.FieldDeleteBy + " = 0")

	if query.ID > 0 {
		tx = tx.Where("id = ?", query.ID)
	}

	if query.Username != "" {
		tx = tx.Where("username = ?", query.Username)
	}

	if query.Nickname != "" {
		tx = tx.Where("nickname ilike ?", "%"+query.Nickname+"%")
	}

	var user domain.Users
	err := tx.Take(&user).Error
	if err != nil {
		return nil, db.ErrFilter(err)
	}

	return &user, nil
}

func (u *UserRepository) FindList(ctx context.Context) ([]*domain.Users, error) {
	var users []*domain.Users
	err := u.GetDB(ctx).WithContext(ctx).Find(&users).Error
	return users, err
}

func (u *UserRepository) FindPage(ctx context.Context, query *domain.UsersPageQuery) (*entity.PaginationResult[domain.Users], error) {
	// 使用通用分页助手
	helper := db.NewPaginationHelper[domain.Users](u.GetDB(ctx))
	result, err := helper.PageQuery(int(query.PageNum), int(query.PageSize), func(dq *gorm.DB) *gorm.DB {
		// 构建基础查询 - 使用 delete_by 条件
		dq = dq.WithContext(ctx).
			Table("users").
			Where(entity.FieldDeleteBy + " = 0")

		// 动态构建查询条件
		if query.Username != "" {
			dq = dq.Where("username LIKE ?", "%"+query.Username+"%")
		}
		if query.Nickname != "" {
			dq = dq.Where("nickname LIKE ?", "%"+query.Nickname+"%")
		}
		if query.Phone != "" {
			dq = dq.Where("phone = ?", query.Phone)
		}
		if query.Email != "" {
			dq = dq.Where("email = ?", query.Email)
		}
		if query.Status != "" {
			dq = dq.Where("status = ?", query.Status)
		}
		if query.UserType != "" {
			dq = dq.Where("type = ?", query.UserType)
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
