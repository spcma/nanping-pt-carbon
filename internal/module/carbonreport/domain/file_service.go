package domain

// fileService 文件领域服务实现
type fileService struct {
	fileRepo FileRepository
}

// NewFileService 创建文件领域服务
func NewFileService(fileRepo FileRepository) FileService {
	return &fileService{
		fileRepo: fileRepo,
	}
}

func (s *fileService) CheckDirExists(path string) (bool, error) {
	_, err := s.fileRepo.GetFile(nil, path)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *fileService) EnsureDirExists(path string, recursive bool) error {
	exists, err := s.CheckDirExists(path)
	if err != nil {
		return err
	}

	if !exists {
		return s.fileRepo.CreateDirectory(nil, path, recursive)
	}

	return nil
}

func (s *fileService) GetFileInfo(path string) (*File, error) {
	return s.fileRepo.GetFile(nil, path)
}

func (s *fileService) ListFiles(path string) ([]File, error) {
	return s.fileRepo.ListDirectory(nil, path)
}
