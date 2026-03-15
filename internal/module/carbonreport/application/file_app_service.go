package application

import (
	"app/internal/module/carbonreport/domain"
	"context"
	"strings"
)

// CheckDirCommand 检查目录命令
type CheckDirCommand struct {
	Path string `json:"path" form:"path" binding:"required"`
}

// CreateDirCommand 创建目录命令
type CreateDirCommand struct {
	Path      string `json:"path" form:"path" binding:"required"`
	Recursive bool   `json:"recursive" form:"recursive"`
}

// ListDirCommand 列出目录命令
type ListDirCommand struct {
	Path string `json:"path" form:"path" query:"path" binding:"required"`
}

// DeleteFileCommand 删除文件命令
type DeleteFileCommand struct {
	Path      string `json:"path" form:"path" query:"path" binding:"required"`
	Recursive bool   `json:"recursive" form:"recursive"`
	Force     bool   `json:"force" form:"force"`
}

// ReadFileCommand 读取文件命令
type ReadFileCommand struct {
	Path string `json:"path" form:"path" query:"path" binding:"required"`
}

// SaveFileCommand 保存文件命令
type SaveFileCommand struct {
	Content  string `json:"content" binding:"required"`
	Dir      string `json:"dir" binding:"required"`
	Filename string `json:"filename" binding:"required"`
}

// UploadFileCommand 上传文件命令
type UploadFileCommand struct {
	Dir      string `form:"dir" binding:"required"`
	Filename string `form:"filename"`
}

// DownloadFileCommand 下载文件命令
type DownloadFileCommand struct {
	Path     string `json:"path" form:"path" query:"path" binding:"required"`
	Filename string `json:"filename" form:"filename" query:"filename"`
}

// FileAppService 文件应用服务
type FileAppService struct {
	fileRepo      domain.FileRepository
	sessionRepo   domain.SessionRepository
	fileService   domain.FileService
	uploadService domain.UploadService
}

// NewFileAppService 创建文件应用服务
func NewFileAppService(
	fileRepo domain.FileRepository,
	sessionRepo domain.SessionRepository,
	fileService domain.FileService,
	uploadService domain.UploadService,
) *FileAppService {
	return &FileAppService{
		fileRepo:      fileRepo,
		sessionRepo:   sessionRepo,
		fileService:   fileService,
		uploadService: uploadService,
	}
}

// CheckDir 检查目录
func (s *FileAppService) CheckDir(ctx context.Context, cmd CheckDirCommand) (*CheckDirResponse, error) {
	exists, err := s.fileService.CheckDirExists(cmd.Path)
	if err != nil {
		return nil, err
	}

	return &CheckDirResponse{
		Path:   cmd.Path,
		Exists: exists,
	}, nil
}

// CreateDir 创建目录
func (s *FileAppService) CreateDir(ctx context.Context, cmd CreateDirCommand) (*CreateDirResponse, error) {
	err := s.fileService.EnsureDirExists(cmd.Path, cmd.Recursive)
	if err != nil {
		return nil, err
	}

	return &CreateDirResponse{
		Path:    cmd.Path,
		Created: true,
	}, nil
}

// ListDir 列出目录
func (s *FileAppService) ListDir(ctx context.Context, cmd ListDirCommand) (*ListDirResponse, error) {
	files, err := s.fileRepo.ListDirectory(ctx, cmd.Path)
	if err != nil {
		return nil, err
	}

	return &ListDirResponse{
		Path:  cmd.Path,
		Files: ConvertFilesToDTOs(files),
	}, nil
}

// DeleteFile 删除文件
func (s *FileAppService) DeleteFile(ctx context.Context, cmd DeleteFileCommand) (*DeleteFileResponse, error) {
	err := s.fileRepo.DeleteFile(ctx, cmd.Path, cmd.Recursive, cmd.Force)
	if err != nil {
		return nil, err
	}

	return &DeleteFileResponse{
		Path:    cmd.Path,
		Deleted: true,
	}, nil
}

// ReadFile 读取文件
func (s *FileAppService) ReadFile(ctx context.Context, cmd ReadFileCommand) (*ReadFileResponse, error) {
	data, size, err := s.fileRepo.ReadFile(ctx, cmd.Path)
	if err != nil {
		return nil, err
	}

	return &ReadFileResponse{
		Path: cmd.Path,
		Size: size,
		Content: func() string {
			if size < 10000 {
				return string(data)
			}
			return "(文件太大，仅显示前 10000 字节)"
		}(),
		Data: data,
	}, nil
}

// SaveFile 保存文件
func (s *FileAppService) SaveFile(ctx context.Context, cmd SaveFileCommand) (*SaveFileResponse, error) {
	ipfsid, err := s.uploadService.UploadFromContent([]byte(cmd.Content), cmd.Dir, cmd.Filename)
	if err != nil {
		return nil, err
	}

	return &SaveFileResponse{
		IPFSID: string(ipfsid),
		Path:   cmd.Dir + "/" + cmd.Filename,
	}, nil
}

// UploadFile 上传文件
func (s *FileAppService) UploadFile(ctx context.Context, localPath string, cmd UploadFileCommand) (*UploadFileResponse, error) {
	filename := cmd.Filename
	if filename == "" {
		filename = cmd.Filename
	}

	ipfsid, err := s.uploadService.UploadFromLocalFile(localPath, cmd.Dir, filename)
	if err != nil {
		return nil, err
	}

	return &UploadFileResponse{
		IPFSID: string(ipfsid),
		Path:   cmd.Dir + "/" + filename,
	}, nil
}

// DownloadFile 下载文件
func (s *FileAppService) DownloadFile(ctx context.Context, cmd DownloadFileCommand) ([]byte, string, error) {
	data, _, err := s.fileRepo.ReadFile(ctx, cmd.Path)
	if err != nil {
		return nil, "", err
	}

	filename := cmd.Filename
	if filename == "" {
		idx := strings.LastIndex(cmd.Path, "/")
		if idx >= 0 {
			filename = cmd.Path[idx+1:]
		} else {
			filename = cmd.Path
		}
	}

	return data, filename, nil
}
