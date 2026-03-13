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

// SysApi API aggregate root
type SysApi struct {
	entity.BaseEntity
	Name       string        `json:"name" gorm:"column:name"`
	Code       string        `json:"code" gorm:"column:code"`
	Uri        string        `json:"uri" gorm:"column:uri"`
	MethodType ApiMethodType `json:"methodType" gorm:"column:method_type"`
	Status     string        `json:"status" gorm:"column:status"`
}

// TableName table name
func (SysApi) TableName() string {
	return "sys_api"
}

// NewSysApi creates a new API
func NewSysApi(name, code, uri, methodType, status string, createUser int64) (*SysApi, error) {
	api := &SysApi{
		BaseEntity: entity.BaseEntity{
			CreateBy:   createUser,
			CreateTime: timeutil.New(),
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
func (a *SysApi) UpdateInfo(name, uri, methodType, status string, userID int64) error {
	a.Name = name
	a.Uri = uri
	a.MethodType = ApiMethodType(methodType)
	a.Status = status
	a.UpdateBy = userID
	a.UpdateTime = timeutil.New()
	return nil
}
