package http

import (
	"app/internal/module/carbonreport/domain"
	"app/internal/module/carbonreport/wire"
	"app/internal/rpc"
	shared_http "app/internal/shared/http"
)

// CreateFsClient 创建文件系统客户端（导出函数）
func CreateFsClient() (*rpc.LApiStub, string, error) {
	return domain.CreateFsClient()
}

// RegisterRoutes 注册 CarbonReport 模块的所有路由
func RegisterRoutes(client *rpc.LApiStub, session string) []shared_http.RouteGroupConfig {
	// 初始化 DDD 组件
	carbonReportDDD := wire.InitFileDDD(client, session)

	// 创建 handlers
	handlers := &Handlers{
		FileHandler: NewFileHandler(carbonReportDDD.AppService),
	}

	return handlers.registerRoutes()
}

// Handlers 包含所有 HTTP 处理器
type Handlers struct {
	FileHandler *FileHandler
}

// registerRoutes 注册 CarbonReport 模块的所有路由（内部方法）
func (h *Handlers) registerRoutes() []shared_http.RouteGroupConfig {
	return []shared_http.RouteGroupConfig{
		// 需要认证的文件管理路由
		{
			Prefix: "/api/v1",
			Handlers: []shared_http.RouteHandler{
				// 目录操作
				{Method: "POST", Path: "/dir/check", Handler: h.FileHandler.CheckDir},
				{Method: "POST", Path: "/dir/create", Handler: h.FileHandler.CreateDir},
				{Method: "GET", Path: "/dir/list", Handler: h.FileHandler.ListDir},
				{Method: "DELETE", Path: "/dir/delete", Handler: h.FileHandler.DeleteFile},
				// 文件操作
				{Method: "GET", Path: "/file/read", Handler: h.FileHandler.ReadFile},
				{Method: "POST", Path: "/file/save", Handler: h.FileHandler.SaveFile},
				{Method: "POST", Path: "/file/upload", Handler: h.FileHandler.UploadFile},
				{Method: "GET", Path: "/file/download", Handler: h.FileHandler.DownloadFile},
				{Method: "DELETE", Path: "/file/delete", Handler: h.FileHandler.DeleteFile},
			},
			AuthType: shared_http.AuthTypeRequired,
		},
	}
}
