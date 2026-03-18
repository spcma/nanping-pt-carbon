package domain

import (
	"app/internal/shared/entity"
	idgen "app/internal/shared/idgen"
	"app/internal/shared/timeutil"
)

type IpfsDetail struct {
	entity.BaseEntity
	DeviceCode     string        `json:"device_code"`
	CollectionTime timeutil.Time `json:"collection_time"`
	TotalDistance  float64       `json:"total_distance"`
	PointCount     int64         `json:"point_count"`
	Filename       string        `json:"filename"`
}

func (*IpfsDetail) TableName() string {
	return "ipfs_detail"
}

// NewIpfsDetail 创建 IPFS 详情
func NewIpfsDetail(deviceCode, filename string, collectionTime timeutil.Time, totalDistance float64, pointCount int64, createUser int64) (*IpfsDetail, error) {
	ipfsDetail := &IpfsDetail{
		BaseEntity: entity.BaseEntity{
			Id:         idgen.NumID(),
			CreateBy:   createUser,
			CreateTime: timeutil.Now(),
		},
		DeviceCode:     deviceCode,
		CollectionTime: collectionTime,
		TotalDistance:  totalDistance,
		PointCount:     pointCount,
		Filename:       filename,
	}
	return ipfsDetail, nil
}

// UpdateInfo 更新 IPFS 详情信息
func (i *IpfsDetail) UpdateInfo(userID int64) error {
	i.UpdateBy = userID
	i.UpdateTime = timeutil.Now()
	return nil
}

// Delete 逻辑删除 IPFS 详情
func (i *IpfsDetail) Delete(userID int64) error {
	i.DeleteBy = userID
	i.DeleteTime = timeutil.Now()
	return nil
}

// IpfsDetailPageQuery IPFS 详情分页查询对象
type IpfsDetailPageQuery struct {
	entity.PaginationQuery
	DeviceCode  string  `json:"device_code" form:"device_code"`
	Filename    string  `json:"filename" form:"filename"`
	MinDistance float64 `json:"min_distance" form:"min_distance"`
	MaxDistance float64 `json:"max_distance" form:"max_distance"`
	SortBy      string  `json:"sortBy" form:"sortBy"`
	SortOrder   string  `json:"sortOrder" form:"sortOrder"` // "asc" or "desc"
}
