package http

import (
	"app/internal/module/ipfs/application"
	shared_http "app/internal/shared/http"

	"github.com/gin-gonic/gin"
)

// 全局服务实例，供定时任务使用
var defaultService *application.Service

// DefaultService 获取默认的碳报告月报服务实例
func DefaultService() *application.Service {
	return defaultService
}

func setDefaultService(service *application.Service) {
	defaultService = service
}

type ipfsRoutes struct {
	ipfsHandler *IpfsHandler
}

func NewIpfsRoutes() shared_http.RouteRegistry {

	service := application.NewService()
	handler := NewIpfsHandler(service)

	setDefaultService(service)

	return &ipfsRoutes{
		ipfsHandler: handler,
	}
}

// RegisterRoutes 注册 IPFS 模块的所有路由
func (i *ipfsRoutes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {

	routerGroup := group.Group("")
	if v, ok := middlewares[shared_http.AuthTypeRequired]; ok {
		routerGroup.Use(v)
	}
	{
		dirRoute := group.Group("dir")
		{
			dirRoute.POST("", i.ipfsHandler.CreateDir)
			dirRoute.DELETE("", i.ipfsHandler.DeleteFile)
			dirRoute.GET("list", i.ipfsHandler.ListDir)
			dirRoute.GET("filestat", i.ipfsHandler.Stat)
			dirRoute.GET("scan", i.ipfsHandler.ScanDir) // 递归扫描目录
		}

		fileRoute := group.Group("file")
		{
			fileRoute.POST("upload", i.ipfsHandler.UploadFile)
			fileRoute.GET("download", i.ipfsHandler.DownloadFile)
			fileRoute.GET("read", i.ipfsHandler.Read)
		}

		calcRoute := group.Group("calc")
		{
			calcRoute.GET("", i.ipfsHandler.CalcDir)
			calcRoute.GET("save", i.ipfsHandler.SaveContentTest)
		}
	}

}
