package infrastructure

import (
	"app/internal/module/ipfs/domain"
	"app/internal/shared/db"
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
	"context"
	"errors"

	"gorm.io/gorm"
)

// IpfsDetailRepository IPFS 详情仓储实现
type IpfsDetailRepository struct {
}

// NewIpfsDetailRepository 创建 IPFS 详情仓储
func NewIpfsDetailRepository(_db *gorm.DB) *IpfsDetailRepository {
	return &IpfsDetailRepository{}
}

// Create 创建 IPFS 详情
func (r *IpfsDetailRepository) Create(ctx context.Context, ipfsDetail *domain.IpfsDetail) error {
	return db.GetDB(ctx).WithContext(ctx).Create(ipfsDetail).Error
}

// Update 更新 IPFS 详情
func (r *IpfsDetailRepository) Update(ctx context.Context, ipfsDetail *domain.IpfsDetail) error {
	return db.GetDB(ctx).WithContext(ctx).Updates(ipfsDetail).Error
}

// UpdateFields 部分更新 IPFS 详情字段
func (r *IpfsDetailRepository) UpdateFields(ctx context.Context, id int64, updates map[string]interface{}) error {
	return db.GetDB(ctx).WithContext(ctx).Model(&domain.IpfsDetail{}).Where("id = ? AND "+entity.FieldDeleteBy+" = 0", id).Updates(updates).Error
}

// Delete 逻辑删除 IPFS 详情
func (r *IpfsDetailRepository) Delete(ctx context.Context, id, uid int64) error {
	updates := map[string]interface{}{
		"deleteBy":   uid,
		"deleteTime": timeutil.Now(),
	}
	return db.GetDB(ctx).WithContext(ctx).Model(&domain.IpfsDetail{}).Where("id = ?", id).Updates(updates).Error
}

// FindByID 根据 ID 查询 IPFS 详情
func (r *IpfsDetailRepository) FindByID(ctx context.Context, id int64) (*domain.IpfsDetail, error) {
	var ipfsDetail domain.IpfsDetail
	err := db.GetDB(ctx).WithContext(ctx).
		Table("ipfs_detail").
		Where("id = ?", id).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&ipfsDetail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ipfsDetail, nil
}

// FindByDeviceCode 根据设备编码查询 IPFS 详情列表
func (r *IpfsDetailRepository) FindByDeviceCode(ctx context.Context, deviceCode string) ([]*domain.IpfsDetail, error) {
	var ipfsDetails []*domain.IpfsDetail
	err := db.GetDB(ctx).WithContext(ctx).
		Table("ipfs_detail").
		Where("device_code = ?", deviceCode).
		Where(entity.FieldDeleteBy + " = 0").
		Find(&ipfsDetails).Error
	return ipfsDetails, err
}

// FindByFilename 根据文件名查询 IPFS 详情
func (r *IpfsDetailRepository) FindByFilename(ctx context.Context, filename string) (*domain.IpfsDetail, error) {
	var ipfsDetail domain.IpfsDetail
	err := db.GetDB(ctx).WithContext(ctx).
		Table("ipfs_detail").
		Where("filename = ?", filename).
		Where(entity.FieldDeleteBy + " = 0").
		Take(&ipfsDetail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &ipfsDetail, nil
}

// FindList 查询 IPFS 详情列表
func (r *IpfsDetailRepository) FindList(ctx context.Context) ([]*domain.IpfsDetail, error) {
	var ipfsDetails []*domain.IpfsDetail
	err := db.GetDB(ctx).WithContext(ctx).
		Table("ipfs_detail").
		Where(entity.FieldDeleteBy + " = 0").
		Find(&ipfsDetails).Error
	return ipfsDetails, err
}

// FindPage 分页查询 IPFS 详情
func (r *IpfsDetailRepository) FindPage(ctx context.Context, query *domain.IpfsDetailPageQuery) (*entity.PaginationResult[*domain.IpfsDetail], error) {
	helper := db.NewPaginationHelper[*domain.IpfsDetail](db.GetDB(ctx))
	result, err := helper.PageQuery(query.PageNum, query.PageSize, func(dq *gorm.DB) *gorm.DB {
		dq = dq.WithContext(ctx).
			Table("ipfs_detail").
			Where(entity.FieldDeleteBy + " = 0")

		// 动态查询条件
		if query.DeviceCode != "" {
			dq = dq.Where("device_code LIKE ?", "%"+query.DeviceCode+"%")
		}
		if query.Filename != "" {
			dq = dq.Where("filename LIKE ?", "%"+query.Filename+"%")
		}
		if query.MinDistance > 0 {
			dq = dq.Where("total_distance >= ?", query.MinDistance)
		}
		if query.MaxDistance > 0 {
			dq = dq.Where("total_distance <= ?", query.MaxDistance)
		}

		return dq
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
