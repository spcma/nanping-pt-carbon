package carbonreportday

import (
	"app/internal/shared/entity"
	idgen "app/internal/shared/idgen"
	"app/internal/shared/timeutil"
)

// CarbonReportDay 碳报告日报聚合根
type CarbonReportDay struct {
	entity.BaseEntity
	Turnover        float64       `json:"turnover" gorm:"column:turnover"`                // 营业额
	Baseline        float64       `json:"baseline" gorm:"column:baseline"`                // 基准值
	CarbonReduction float64       `json:"carbonReduction" gorm:"column:carbon_reduction"` // 碳减排量
	CollectionDate  timeutil.Time `json:"collection_date" gorm:"column:collection_date"`  // 数据采集日期
}

// TableName 表名
func (*CarbonReportDay) TableName() string {
	return "carbon_report_day"
}

// NewCarbonReportDay 创建新的碳报告日报
func NewCarbonReportDay(turnover, baseline float64, collectionDate timeutil.Time, createUser int64) (*CarbonReportDay, error) {
	report := &CarbonReportDay{
		BaseEntity: entity.BaseEntity{
			Id:         idgen.NumID(),
			CreateBy:   createUser,
			CreateTime: timeutil.Now(),
		},
		Turnover:       turnover,
		Baseline:       baseline,
		CollectionDate: collectionDate,
	}
	return report, nil
}

// UpdateInfo 更新碳报告日报信息
func (r *CarbonReportDay) UpdateInfo(userID int64) error {
	r.UpdateBy = userID
	r.UpdateTime = timeutil.Now()
	return nil
}

// CarbonReportDayPageQuery 碳报告日报分页查询对象
type CarbonReportDayPageQuery struct {
	entity.PaginationQuery
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"` // "asc" or "desc"
	StartDate string `json:"startDate"` // 开始日期
	EndDate   string `json:"endDate"`   // 结束日期
}
