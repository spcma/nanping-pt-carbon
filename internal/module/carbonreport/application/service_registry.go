package application

import (
	"app/internal/module/carbonreport/domain"
	"app/internal/shared/event"
)

// ServiceRegistry 服务注册表
type ServiceRegistry struct {
	fileService   domain.FileService
	uploadService domain.UploadService
	eventBus      event.EventBus
}

// GetFileService 获取文件服务
func (r *ServiceRegistry) GetFileService() domain.FileService {
	return r.fileService
}

// GetUploadService 获取上传服务
func (r *ServiceRegistry) GetUploadService() domain.UploadService {
	return r.uploadService
}

// GetEventBus 获取事件总线
func (r *ServiceRegistry) GetEventBus() event.EventBus {
	return r.eventBus
}
