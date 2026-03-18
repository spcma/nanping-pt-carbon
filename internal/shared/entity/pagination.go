package entity

type PaginationQuery struct {
	PageNum  int `json:"pageNum" form:"pageNum"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

func (p *PaginationQuery) Fixed() {
	if p.PageNum <= 0 {
		p.PageNum = 1
	}

	if p.PageSize <= 0 {
		p.PageSize = 10
	}
}

// PaginationResult 分页查询结果
type PaginationResult[T any] struct {
	Data       []T   `json:"data"`       // 数据列表
	Total      int64 `json:"total"`      // 总记录数
	PageNum    int   `json:"pageNum"`    // 当前页码
	PageSize   int   `json:"pageSize"`   // 每页大小
	TotalPages int   `json:"totalPages"` // 总页数
}

// NewPaginationResult 创建分页结果
func NewPaginationResult[T any](data []T, total int64, pageNum, pageSize int) *PaginationResult[T] {
	totalPages := 0
	if pageSize > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return &PaginationResult[T]{
		Data:       data,
		Total:      total,
		PageNum:    pageNum,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}

// HasNext 是否有下一页
func (p *PaginationResult[T]) HasNext() bool {
	return p.PageNum < p.TotalPages
}

// HasPrev 是否有上一页
func (p *PaginationResult[T]) HasPrev() bool {
	return p.PageNum > 1
}

// IsFirst 是否第一页
func (p *PaginationResult[T]) IsFirst() bool {
	return p.PageNum == 1
}

// IsLast 是否最后一页
func (p *PaginationResult[T]) IsLast() bool {
	return p.PageNum == p.TotalPages || p.TotalPages == 0
}
