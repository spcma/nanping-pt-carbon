package application

import (
	"app/internal/module/carbonreport/domain"
	"context"
)

// FileApplicationService 文件应用服务
type FileApplicationService struct {
	fileRepo    domain.FileRepository
	sessionRepo domain.SessionRepository
}

// NewFileApplicationService 创建文件应用服务
func NewFileApplicationService(fileRepo domain.FileRepository, sessionRepo domain.SessionRepository) *FileApplicationService {
	return &FileApplicationService{
		fileRepo:    fileRepo,
		sessionRepo: sessionRepo,
	}
}

// CheckDirExists 检查目录是否存在
func (s *FileApplicationService) CheckDirExists(ctx context.Context, path string) (bool, error) {
	_, err := s.fileRepo.GetFile(ctx, path)
	if err != nil {
		if err == domain.ErrDirNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateDirectory 创建目录
func (s *FileApplicationService) CreateDirectory(ctx context.Context, path string, recursive bool) error {
	return s.fileRepo.CreateDirectory(ctx, path, recursive)
}

// ListDirectory 列出目录
func (s *FileApplicationService) ListDirectory(ctx context.Context, path string) ([]FileDTO, error) {
	files, err := s.fileRepo.ListDirectory(ctx, path)
	if err != nil {
		return nil, err
	}
	return ConvertFilesToDTOs(files), nil
}

// DeleteFile 删除文件
func (s *FileApplicationService) DeleteFile(ctx context.Context, path string, recursive, force bool) error {
	return s.fileRepo.DeleteFile(ctx, path, recursive, force)
}

// ReadFile 读取文件
func (s *FileApplicationService) ReadFile(ctx context.Context, filePath string) ([]byte, int64, error) {
	return s.fileRepo.ReadFile(ctx, filePath)
}

// SaveContent 保存内容到文件
func (s *FileApplicationService) SaveContent(ctx context.Context, content, dir, filename string) (string, error) {
	ipfsid, err := s.fileRepo.SaveContent(ctx, []byte(content), dir, filename)
	if err != nil {
		return "", err
	}
	return string(ipfsid), nil
}

// SaveLocalFile 保存本地文件到 NPFS
func (s *FileApplicationService) SaveLocalFile(ctx context.Context, localPath, fsDir, filename string) (string, error) {
	ipfsid, err := s.fileRepo.SaveLocalFile(ctx, localPath, fsDir, filename)
	if err != nil {
		return "", err
	}
	return string(ipfsid), nil
}
