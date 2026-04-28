package carbonreportdetail

import (
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
)

type CarbonReportDetail struct {
	entity.BaseEntity
	TraceCode      string        `json:"trace_code" gorm:"column:trace_code"`
	Mileage        float64       `json:"mileage" gorm:"column:mileage"`
	Passenger      int64         `json:"passenger" gorm:"column:passenger"`
	Turnover       float64       `json:"turnover" gorm:"column:turnover"`
	CollectionTime timeutil.Time `json:"collection_time" gorm:"column:collection_time"`
	Hash           string        `json:"hash" gorm:"column:hash"`
}

func (*CarbonReportDetail) TableName() string {
	return "carbon_report_detail"
}
