package domain

import (
	"app/internal/shared/entity"
)

// CarbonReportDayPageQuery 碳报告日报分页查询对象
type CarbonReportDayPageQuery struct {
	entity.PaginationQuery
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"` // "asc" or "desc"
	StartDate string `json:"startDate"` // 开始日期
	EndDate   string `json:"endDate"`   // 结束日期
}
