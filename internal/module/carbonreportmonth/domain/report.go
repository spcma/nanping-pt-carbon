package domain

import (
	"app/internal/shared/entity"
	idgen "app/internal/shared/idgen"
	"app/internal/shared/timeutil"
)

// CarbonReportMonth 碳月报聚合根
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

// NewCarbonReportMonth 创建新的碳月报（工厂方法）
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
		CarbonReduction:   0, // 初始为0，后续计算
	}
	return report, nil
}

// UpdateInfo 更新碳月报信息（领域行为）
func (r *CarbonReportMonth) UpdateInfo(userID int64) error {
	r.UpdateBy = userID
	r.UpdateTime = timeutil.Now()
	return nil
}

// CalculateCarbonReduction 计算碳减排量（领域行为）
func (r *CarbonReportMonth) CalculateCarbonReduction() {
	// 碳减排量 = 基准值 - 能耗
	r.CarbonReduction = r.Baseline - r.EnergyConsumption
}

// SetEnergyConsumption 设置能耗并重新计算碳减排量
func (r *CarbonReportMonth) SetEnergyConsumption(energyConsumption float64, userID int64) {
	r.EnergyConsumption = energyConsumption
	r.CalculateCarbonReduction()
	r.UpdateInfo(userID)
}
