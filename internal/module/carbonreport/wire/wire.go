package wire

import (
	"app/internal/module/carbonreport/application"
	"app/internal/module/carbonreport/domain"
	"app/internal/module/carbonreport/infrastructure"
	"app/internal/module/ipfs/rpc"
)

type FileDDD struct {
	FileRepo      domain.FileRepository
	SessionRepo   domain.SessionRepository
	FileService   domain.FileService
	UploadService domain.UploadService
	AppService    *application.FileAppService
}

// InitFileDDD initializes file DDD components
func InitFileDDD(client *rpc.LApiStub, session string) *FileDDD {
	// 1. 初始化仓储层
	fileRepo := infrastructure.NewFileRepository(client, session)
	sessionRepo := infrastructure.NewSessionRepository(client, session)

	// 2. 初始化领域服务层
	fileService := domain.NewFileService(fileRepo)
	uploadService := domain.NewUploadService(fileRepo)

	// 3. 初始化应用服务层
	appService := application.NewFileAppService(fileRepo, sessionRepo, fileService, uploadService)

	return &FileDDD{
		FileRepo:      fileRepo,
		SessionRepo:   sessionRepo,
		FileService:   fileService,
		UploadService: uploadService,
		AppService:    appService,
	}
}
