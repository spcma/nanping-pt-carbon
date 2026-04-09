package application

import (
	"app/internal/module/ipfs/domain"
	"app/internal/shared/entity"
	"context"
)

// CreateIpfsDetailCommand 创建 IPFS 详情命令
type CreateIpfsDetailCommand struct {
	DeviceCode     string  `json:"device_code"`
	Filename       string  `json:"filename"`
	CollectionTime string  `json:"collection_time"`
	TotalDistance  float64 `json:"total_distance"`
	PointCount     int64   `json:"point_count"`
	UserID         int64   `json:"userId"`
}

// UpdateIpfsDetailCommand 更新 IPFS 详情命令
type UpdateIpfsDetailCommand struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"userId"`
	// 可根据需要添加其他可更新字段
}

// IpfsDetailAppService IPFS 详情应用服务
type IpfsDetailAppService struct {
	repo IpfsDetailRepo
}

// NewIpfsDetailAppService 创建 IPFS 详情应用服务
func NewIpfsDetailAppService(repo IpfsDetailRepo) *IpfsDetailAppService {
	return &IpfsDetailAppService{repo: repo}
}

// CreateIpfsDetail 创建 IPFS 详情
func (s *IpfsDetailAppService) CreateIpfsDetail(ctx context.Context, cmd CreateIpfsDetailCommand) (int64, error) {
	//collectionTime, err := timeutil.FromString(cmd.CollectionTime)
	//if err != nil {
	//	return 0, err
	//}

	//ipfsDetail, err := domain.NewIpfsDetail(
	//	cmd.DeviceCode,
	//	cmd.Filename,
	//	collectionTime,
	//	cmd.TotalDistance,
	//	cmd.PointCount,
	//	cmd.UserID,
	//)
	//if err != nil {
	//	return 0, err
	//}
	//
	//err = s.repo.Create(ctx, ipfsDetail)
	//if err != nil {
	//	return 0, err
	//}
	//
	//return ipfsDetail.Id, nil

	return 0, nil
}

// UpdateIpfsDetail 更新 IPFS 详情
func (s *IpfsDetailAppService) UpdateIpfsDetail(ctx context.Context, cmd UpdateIpfsDetailCommand) error {
	ipfsDetail, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}
	if ipfsDetail == nil {
		return domain.ErrIpfsDetailNotFound
	}
	return ipfsDetail.UpdateInfo(cmd.UserID)
}

// DeleteIpfsDetail 删除 IPFS 详情（逻辑删除）
func (s *IpfsDetailAppService) DeleteIpfsDetail(ctx context.Context, id int64, userID int64) error {
	ipfsDetail, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if ipfsDetail == nil {
		return domain.ErrIpfsDetailNotFound
	}
	return ipfsDetail.Delete(userID)
}

// GetIpfsDetailByID 根据 ID 获取 IPFS 详情
func (s *IpfsDetailAppService) GetIpfsDetailByID(ctx context.Context, id int64) (*domain.IpfsDetail, error) {
	return s.repo.FindByID(ctx, id)
}

// GetIpfsDetailByDeviceCode 根据设备编码获取 IPFS 详情列表
func (s *IpfsDetailAppService) GetIpfsDetailByDeviceCode(ctx context.Context, deviceCode string) ([]*domain.IpfsDetail, error) {
	return s.repo.FindByDeviceCode(ctx, deviceCode)
}

// GetIpfsDetailByFilename 根据文件名获取 IPFS 详情
func (s *IpfsDetailAppService) GetIpfsDetailByFilename(ctx context.Context, filename string) (*domain.IpfsDetail, error) {
	return s.repo.FindByFilename(ctx, filename)
}

// GetIpfsDetailPage 分页查询 IPFS 详情
func (s *IpfsDetailAppService) GetIpfsDetailPage(ctx context.Context, query *domain.IpfsDetailPageQuery) (*entity.PaginationResult[domain.IpfsDetail], error) {
	return s.repo.FindPage(ctx, query)
}
