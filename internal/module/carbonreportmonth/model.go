package carbonreportmonth

import (
	"app/internal/shared/entity"
	idgen "app/internal/shared/idgen"
	"app/internal/shared/timeutil"
)

type CarbonReportMonth struct {
	entity.BaseEntity
	CollectionDate    timeutil.Time `json:"collection_date" gorm:"column:collection_date"`      // 数据采集日期
	Turnover          float64       `json:"turnover" gorm:"column:turnover"`                    // 周转量
	Baseline          float64       `json:"baseline" gorm:"column:baseline"`                    // 基准值
	EnergyConsumption float64       `json:"energyConsumption" gorm:"column:energy_consumption"` // 能耗, 人工填入
	CarbonReduction   float64       `json:"carbonReduction" gorm:"column:carbon_reduction"`     // 碳减排量
}

// TableName 表名
func (*CarbonReportMonth) TableName() string {
	return "carbon_report_month"
}

// NewCarbonReportMonth 创建新的碳报告月报
func NewCarbonReportMonth(turnover, baseline, energyConsumption float64, collectionDate timeutil.Time, createUser int64) (*CarbonReportMonth, error) {
	report := &CarbonReportMonth{
		BaseEntity: entity.BaseEntity{
			Id:         idgen.NumID(),
			CreateBy:   createUser,
			CreateTime: timeutil.Now(),
		},
		Turnover:          turnover,
		Baseline:          baseline,
		EnergyConsumption: energyConsumption,
		CollectionDate:    collectionDate,
	}
	return report, nil
}

// UpdateInfo 更新碳报告月报信息
func (r *CarbonReportMonth) UpdateInfo(userID int64) error {
	r.UpdateBy = userID
	r.UpdateTime = timeutil.Now()
	return nil
}

// CarbonReportMonthPageQuery 碳报告月报分页查询对象
type CarbonReportMonthPageQuery struct {
	entity.PaginationQuery
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"` // "asc" or "desc"
	StartDate string `json:"startDate"` // 开始日期
	EndDate   string `json:"endDate"`   // 结束日期
}
