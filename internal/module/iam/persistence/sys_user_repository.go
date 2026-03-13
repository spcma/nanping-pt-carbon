package persistence

import (
	"app/internal/module/iam/domain"
	"app/internal/shared/db"
	"app/internal/shared/timeutil"
	"context"
	"strings"

	"gorm.io/gorm"
)

type sysUserRepository struct {
	db *gorm.DB
}

// NewSysUserRepository creates user repository
func NewSysUserRepository(db *gorm.DB) domain.SysUserRepository {
	return &sysUserRepository{db: db}
}

func (r *sysUserRepository) Create(ctx context.Context, user *domain.SysUser) error {
	return db.WithContext(ctx, r.db).Create(user).Error
}

func (r *sysUserRepository) Update(ctx context.Context, user *domain.SysUser) error {
	return db.WithContext(ctx, r.db).Save(user).Error
}

func (r *sysUserRepository) Delete(ctx context.Context, id int64, userID int64) error {
	// 逻辑删除：只更新 delete_user 和 delete_time 字段
	updates := map[string]interface{}{
		"delete_by":   userID,
		"delete_time": timeutil.New(),
	}
	return db.WithContext(ctx, r.db).Model(&domain.SysUser{}).Where("id = ?", id).Updates(updates).Error
}

func (r *sysUserRepository) FindByID(ctx context.Context, id int64) (*domain.SysUser, error) {
	var user domain.SysUser
	err := db.WithContext(ctx, r.db).Where("id = ? AND delete_time IS NULL", id).First(&user).Error
	return &user, err
}

func (r *sysUserRepository) FindByUsername(ctx context.Context, username string) (*domain.SysUser, error) {
	var user domain.SysUser
	err := db.WithContext(ctx, r.db).Where("username = ? AND delete_time IS NULL", username).First(&user).Error
	return &user, err
}

func (r *sysUserRepository) FindPage(ctx context.Context, query *domain.SysUserPageQuery) ([]*domain.SysUser, int64, error) {
	var users []*domain.SysUser
	var total int64

	tx := db.WithContext(ctx, r.db).Model(&domain.SysUser{}).Where("delete_time IS NULL")

	// 动态构建查询条件
	if query.Username != "" {
		tx = tx.Where("username LIKE ?", "%"+query.Username+"%")
	}
	if query.Nickname != "" {
		tx = tx.Where("nickname LIKE ?", "%"+query.Nickname+"%")
	}
	if query.Phone != "" {
		tx = tx.Where("phone = ?", query.Phone)
	}
	if query.Email != "" {
		tx = tx.Where("email = ?", query.Email)
	}
	if query.Status != "" {
		tx = tx.Where("status = ?", query.Status)
	}
	if query.UserType != "" {
		tx = tx.Where("user_type = ?", query.UserType)
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
	err := tx.Offset(int(offset)).Limit(int(query.PageSize)).Find(&users).Error
	return users, total, err
}
