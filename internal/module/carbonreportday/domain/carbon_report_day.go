package domain

import (
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
)

// CarbonReportDay 碳报告日报聚合根
type CarbonReportDay struct {
	entity.BaseEntity
	// TODO: 根据实际业务需求添加业务字段
	// 例如：
	// ReportDate timeutil.Time `json:"report_date" gorm:"column:report_date"` // 报告日期
	// TotalEmission float64 `json:"total_emission" gorm:"column:total_emission"` // 总排放量
}

// TableName 表名
func (*CarbonReportDay) TableName() string {
	return "carbon_report_day"
}

// NewCarbonReportDay 创建新的碳报告日报
func NewCarbonReportDay(createUser int64) (*CarbonReportDay, error) {
	report := &CarbonReportDay{
		BaseEntity: entity.BaseEntity{
			CreateBy:   createUser,
			CreateTime: timeutil.Now(),
		},
	}
	return report, nil
}

// UpdateInfo 更新碳报告日报信息
func (r *CarbonReportDay) UpdateInfo(userID int64) error {
	r.UpdateBy = userID
	r.UpdateTime = timeutil.Now()
	return nil
}
