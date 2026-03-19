package http

import (
	"app/internal/config"
	"app/internal/module/ipfs/application"
	shared_http "app/internal/shared/http"
	"app/internal/shared/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ipfsRoutes struct {
	ipfsHandler *IpfsHandler
}

func NewIpfsRoutes(db *gorm.DB, remoteDB *gorm.DB) shared_http.RouteRegistry {

	var appService *application.Service
	//	是否开启ipfs本地服务
	if config.GlobalConfig.Ipfs.Status {
		fsClient, sessionId, err := application.CreateFsClient()
		if err != nil {
			logger.Error("http", "Failed to create IPFS client: "+err.Error())
			panic(err)
		}
		appService = application.NewService(db, remoteDB, fsClient, sessionId)
	} else {
		appService = application.NewService(db, remoteDB, nil, "")
	}

	ipfsHandler, err := NewIpfsHandler(appService)
	if err != nil {
		logger.Error("http", "Failed to create IPFS handler: "+err.Error())
		panic(err)
	}

	return &ipfsRoutes{
		ipfsHandler: ipfsHandler,
	}
}

// RegisterRoutes 注册 IPFS 模块的所有路由
func (i *ipfsRoutes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {

	dirRoute := group.Group("dir")
	if v, ok := middlewares[shared_http.AuthTypeRequired]; ok {
		dirRoute.Use(v)
	}
	{
		dirRoute.POST("", i.ipfsHandler.CreateDir)
		dirRoute.DELETE("", i.ipfsHandler.DeleteFile)
		dirRoute.GET("list", i.ipfsHandler.ListDir)
		dirRoute.POST("upload", i.ipfsHandler.UploadFileHandler)
		dirRoute.GET("download", i.ipfsHandler.DownloadFileHandler)

		dirRoute.POST("handle", i.service.HandleWithDir)
		dirRoute.GET("hhh", i.service.HHH)
		dirRoute.GET("h1", i.service.H1)
	}

	fileRoute := group.Group("file")
	if v, ok := middlewares[shared_http.AuthTypeRequired]; ok {
		fileRoute.Use(v)
	}
	{
		fileRoute.GET("read", i.service.ReadFile)
		fileRoute.POST("save", i.service.SaveFile)
		fileRoute.POST("upload", i.service.UploadFile)
		fileRoute.GET("download", i.service.DownloadFile)
		fileRoute.DELETE("delete", i.service.DeleteFile)
	}
}
