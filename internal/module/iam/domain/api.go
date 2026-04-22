package domain

import (
	"app/internal/shared/entity"
	"app/internal/shared/timeutil"
)

// ApiMethodType API method type
type ApiMethodType string

const (
	MethodGet    ApiMethodType = "GET"    // GET
	MethodPost   ApiMethodType = "POST"   // POST
	MethodPut    ApiMethodType = "PUT"    // PUT
	MethodDelete ApiMethodType = "DELETE" // DELETE
)

// Api API aggregate root
type Api struct {
	entity.BaseEntity
	Name       string        `json:"name" gorm:"column:name"`
	Code       string        `json:"code" gorm:"column:code"`
	Uri        string        `json:"uri" gorm:"column:uri"`
	MethodType ApiMethodType `json:"methodType" gorm:"column:method_type"`
	Status     string        `json:"status" gorm:"column:status"`
}

// TableName table name
func (Api) TableName() string {
	return "sys_api"
}

// NewApi creates a new API
func NewApi(name, code, uri, methodType, status string, createUser int64) (*Api, error) {
	api := &Api{
		BaseEntity: entity.BaseEntity{
			CreateBy:   createUser,
			CreateTime: timeutil.Now(),
		},
		Name:       name,
		Code:       code,
		Uri:        uri,
		MethodType: ApiMethodType(methodType),
		Status:     status,
	}
	return api, nil
}

// UpdateInfo updates API info
func (a *Api) UpdateInfo(name, uri, methodType, status string, userID int64) error {
	a.Name = name
	a.Uri = uri
	a.MethodType = ApiMethodType(methodType)
	a.Status = status
	a.UpdateBy = userID
	a.UpdateTime = timeutil.Now()
	return nil
}

// SysApiPageQuery system API page query object
type SysApiPageQuery struct {
	PageNum    int64  `json:"pageNum" binding:"required,min=1"`
	PageSize   int64  `json:"pageSize" binding:"required,min=1,max=100"`
	Name       string `json:"name"`
	Code       string `json:"code"`
	Uri        string `json:"uri"`
	MethodType string `json:"methodType"`
	Status     string `json:"status"`
	SortBy     string `json:"sortBy"`
	SortOrder  string `json:"sortOrder" binding:"oneof=asc desc"` // "asc" or "desc"
}
