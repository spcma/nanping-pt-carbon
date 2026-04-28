package domain

import (
	"app/internal/shared/entity"
	idgen "app/internal/shared/idgen"
	"app/internal/shared/timeutil"
)

// CarbonReportDay 碳日报聚合根
type CarbonReportDay struct {
	entity.BaseEntity
	TraceCode       string        `json:"trace_code" gorm:"column:trace_code"` // 追溯码
	Mileage         float64       `json:"mileage" gorm:"column:mileage"`
	Turnover        float64       `json:"turnover" gorm:"column:turnover"`               // 营业额
	CollectionDate  timeutil.Time `json:"collection_date" gorm:"column:collection_date"` // 数据采集日期
	Hash            string        `json:"hash" gorm:"column:hash"`
	Baseline        float64       `json:"baseline" gorm:"column:baseline"`                // 基准值
	CarbonReduction float64       `json:"carbonReduction" gorm:"column:carbon_reduction"` // 碳减排量
}

// TableName 表名
func (*CarbonReportDay) TableName() string {
	return "carbon_report_day"
}

// NewCarbonReportDay 创建新的碳日报（工厂方法）
func NewCarbonReportDay(turnover, baseline float64, collectionDate timeutil.Time, createUser int64) (*CarbonReportDay, error) {
	report := &CarbonReportDay{
		BaseEntity: entity.BaseEntity{
			Id:         idgen.NumID(),
			CreateBy:   createUser,
			CreateTime: timeutil.Now(),
		},
		Turnover:        turnover,
		Baseline:        baseline,
		CollectionDate:  collectionDate,
		CarbonReduction: 0, // 初始为0，后续计算
	}
	return report, nil
}

// UpdateInfo 更新碳日报信息（领域行为）
func (r *CarbonReportDay) UpdateInfo(userID int64) error {
	r.UpdateBy = userID
	r.UpdateTime = timeutil.Now()
	return nil
}

// CalculateCarbonReduction 计算碳减排量（领域行为）
func (r *CarbonReportDay) CalculateCarbonReduction() {
	// 碳减排量 = 基准值 - 实际排放（这里简化处理）
	r.CarbonReduction = r.Baseline - r.Turnover*0.1 // 示例公式
}

// SetHash 设置哈希值
func (r *CarbonReportDay) SetHash(hash string) {
	r.Hash = hash
}
