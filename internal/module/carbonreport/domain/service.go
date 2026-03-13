package domain

// FileService 文件领域服务
type FileService interface {
	// CheckDirExists 检查目录是否存在
	CheckDirExists(path string) (bool, error)

	// EnsureDirExists 确保目录存在
	EnsureDirExists(path string, recursive bool) error

	// GetFileInfo 获取文件信息
	GetFileInfo(path string) (*File, error)

	// ListFiles 列出文件列表
	ListFiles(path string) ([]File, error)
}

// UploadService 上传领域服务
type UploadService interface {
	// UploadFromContent 从内容上传
	UploadFromContent(content []byte, dirPath, filename string) (IPFSID, error)

	// UploadFromLocalFile 从本地文件上传
	UploadFromLocalFile(localPath, fsDir, filename string) (IPFSID, error)
}
