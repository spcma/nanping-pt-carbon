package http

import (
	"app/internal/module/ipfs/application"
	"app/internal/shared/db"
	shared_http "app/internal/shared/http"
	"app/internal/shared/logger"

	"github.com/gin-gonic/gin"
)

type ipfsRoutes struct {
	ipfsHandler *IpfsHandler
}

func NewIpfsRoutes() shared_http.RouteRegistry {
	dbt := db.Default()
	remoteDB := db.Remote()

	var appService *application.Service

	appService = application.NewService(dbt, remoteDB)
	if appService == nil {
		logger.Error("http", "Failed to create IPFS service")
		panic("failed to create IPFS service")
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
			dirRoute.GET("check", i.ipfsHandler.CheckDir)
			dirRoute.GET("scan", i.ipfsHandler.ScanDir) // 递归扫描目录
		}

		fileRoute := group.Group("file")
		{
			fileRoute.POST("upload", i.ipfsHandler.UploadFile)
			fileRoute.GET("download", i.ipfsHandler.DownloadFile)
			fileRoute.GET("read", i.ipfsHandler.ReadIpfs)
		}

		calcRoute := group.Group("calc")
		{
			calcRoute.GET("", i.ipfsHandler.CalcDir)
			calcRoute.GET("save", i.ipfsHandler.SaveContentTest)
		}
	}

}
