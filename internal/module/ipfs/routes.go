package ipfs

import (
	"app/internal/config"
	shared_http "app/internal/shared/http"
	"app/internal/shared/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ipfsRoutes struct {
	service *Service
}

func NewIpfsRoutes(db *gorm.DB) shared_http.RouteRegistry {
	newService := NewService(db, nil, "")

	if config.GlobalConfig.Ipfs.Status {
		fsClient, sessionId, err := CreateFsClient()
		if err != nil {
			logger.Error("http", "Failed to create IPFS client: "+err.Error())
			panic(err)
		}

		newService.client = fsClient
		newService.session = sessionId
	}

	return &ipfsRoutes{
		service: newService,
	}
}

// RegisterRoutes 注册 IPFS 模块的所有路由
func (i *ipfsRoutes) RegisterRoutes(group *gin.RouterGroup, middlewares map[shared_http.AuthType]gin.HandlerFunc) {

	dirRoute := group.Group("dir")
	if v, ok := middlewares[shared_http.AuthTypeRequired]; ok {
		dirRoute.Use(v)
	}
	{
		dirRoute.POST("check", i.service.CheckDir)
		dirRoute.POST("create", i.service.CreateDir)
		dirRoute.GET("list", i.service.ListDir)
		dirRoute.POST("handle", i.service.HandleWithDir)
		dirRoute.DELETE("delete", i.service.DeleteFile)
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
