package domain

import "context"

// CarbonReportDayPageQuery 碳报告日报分页查询对象
type CarbonReportDayPageQuery struct {
	PageNum   int64  `json:"pageNum" binding:"required,min=1"`
	PageSize  int64  `json:"pageSize" binding:"required,min=1,max=100"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder" binding:"oneof=asc desc"` // "asc" or "desc"
	// TODO: 根据实际业务需求添加查询条件
	// 例如：
	// StartDate string `json:"startDate"` // 开始日期
	// EndDate   string `json:"endDate"`   // 结束日期
}

// CarbonReportDayRepository 碳报告日报仓储接口
type CarbonReportDayRepository interface {
	Create(ctx context.Context, report *CarbonReportDay) error
	Update(ctx context.Context, report *CarbonReportDay) error
	Delete(ctx context.Context, id int64, userID int64) error
	FindByID(ctx context.Context, id int64) (*CarbonReportDay, error)
	FindPage(ctx context.Context, query CarbonReportDayPageQuery) ([]*CarbonReportDay, int64, error)
	// TODO: 根据实际业务需求添加其他查询方法
	// FindByDate(ctx context.Context, date timeutil.Time) (*CarbonReportDay, error)
}
