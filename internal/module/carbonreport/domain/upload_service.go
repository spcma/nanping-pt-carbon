package domain

// uploadService 上传领域服务实现
type uploadService struct {
	fileRepo FileRepository
}

// NewUploadService 创建上传领域服务
func NewUploadService(fileRepo FileRepository) UploadService {
	return &uploadService{
		fileRepo: fileRepo,
	}
}

func (s *uploadService) UploadFromContent(content []byte, dirPath, filename string) (IPFSID, error) {
	return s.fileRepo.SaveContent(nil, content, dirPath, filename)
}

func (s *uploadService) UploadFromLocalFile(localPath, fsDir, filename string) (IPFSID, error) {
	return s.fileRepo.SaveLocalFile(nil, localPath, fsDir, filename)
}
