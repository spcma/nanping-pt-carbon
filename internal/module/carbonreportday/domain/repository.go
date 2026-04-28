package domain

import (
	"app/internal/shared/entity"
	"context"
)

// CarbonReportDayRepository 碳日报仓储接口（领域层定义）
type CarbonReportDayRepository interface {
	Create(ctx context.Context, report *CarbonReportDay) error
	Update(ctx context.Context, report *CarbonReportDay) error
	Delete(ctx context.Context, id int64, userID int64) error

	FindByID(ctx context.Context, id int64) (*CarbonReportDay, error)
	FindList(ctx context.Context) ([]*CarbonReportDay, error)
	FindPage(ctx context.Context, query *CarbonReportDayPageQuery) (*entity.PaginationResult[*CarbonReportDay], error)

	// FindByMonth 根据年月查询该月的所有日报数据
	FindByMonth(ctx context.Context, year int, month int) ([]*CarbonReportDay, error)
}

// CarbonReportDayPageQuery 碳日报分页查询对象
type CarbonReportDayPageQuery struct {
	entity.PaginationQuery
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"` // "asc" or "desc"
	StartDate string `json:"startDate"` // 开始日期
	EndDate   string `json:"endDate"`   // 结束日期
}
