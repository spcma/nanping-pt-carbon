package domain

import (
	"app/internal/shared/entity"
	"context"
)

// CarbonReportMonthRepository 碳月报仓储接口（领域层定义）
type CarbonReportMonthRepository interface {
	Create(ctx context.Context, report *CarbonReportMonth) error
	Update(ctx context.Context, report *CarbonReportMonth) error
	Delete(ctx context.Context, id int64, userID int64) error

	FindByID(ctx context.Context, id int64) (*CarbonReportMonth, error)
	FindList(ctx context.Context) ([]*CarbonReportMonth, error)
	FindPage(ctx context.Context, query *CarbonReportMonthPageQuery) (*entity.PaginationResult[*CarbonReportMonth], error)

	// FindByMonth 根据年月查询月报
	FindByMonth(ctx context.Context, year int, month int) (*CarbonReportMonth, error)
}

// CarbonReportMonthPageQuery 碳月报分页查询对象
type CarbonReportMonthPageQuery struct {
	entity.PaginationQuery
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"` // "asc" or "desc"
	StartDate string `json:"startDate"` // 开始日期
	EndDate   string `json:"endDate"`   // 结束日期
}
