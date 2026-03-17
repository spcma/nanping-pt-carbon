package ipfs

import (
	"app/internal/module/ipfs/rpc"
	shared_http "app/internal/shared/http"
)

// RegisterRoutes 注册 IPFS 模块的所有路由
func RegisterRoutes(client *rpc.LApiStub, session string) []shared_http.RouteGroupConfig {
	// 创建服务实例
	service := NewService(client, session)

	return []shared_http.RouteGroupConfig{
		{
			Prefix: "",
			Handlers: []shared_http.RouteHandler{
				// 目录操作
				{Method: "POST", Path: "/dir/check", Handler: service.CheckDir},
				{Method: "POST", Path: "/dir/create", Handler: service.CreateDir},
				{Method: "GET", Path: "/dir/list", Handler: service.ListDir},
				{Method: "GET", Path: "/dir/handle", Handler: service.HandleWithDir},
				{Method: "DELETE", Path: "/dir/delete", Handler: service.DeleteFile},
				// 文件操作
				{Method: "GET", Path: "/file/read", Handler: service.ReadFile},
				{Method: "POST", Path: "/file/save", Handler: service.SaveFile},
				{Method: "POST", Path: "/file/upload", Handler: service.UploadFile},
				{Method: "GET", Path: "/file/download", Handler: service.DownloadFile},
				{Method: "DELETE", Path: "/file/delete", Handler: service.DeleteFile},
			},
			AuthType: shared_http.AuthTypeRequired,
		},
	}
}
