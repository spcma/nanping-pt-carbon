package domain

import "context"

// FileRepository 文件仓储接口
type FileRepository interface {
	// GetFile 获取文件
	GetFile(ctx context.Context, path string) (*File, error)

	// ListDirectory 列出目录
	ListDirectory(ctx context.Context, path string) ([]File, error)

	// CreateDirectory 创建目录
	CreateDirectory(ctx context.Context, path string, recursive bool) error

	// DeleteFile 删除文件
	DeleteFile(ctx context.Context, path string, recursive, force bool) error

	// SaveContent 保存内容到文件
	SaveContent(ctx context.Context, content []byte, dirPath, filename string) (IPFSID, error)

	// SaveLocalFile 保存本地文件到 NPFS
	SaveLocalFile(ctx context.Context, localPath, fsDir, filename string) (IPFSID, error)

	// ReadFile 读取文件内容
	ReadFile(ctx context.Context, filePath string) ([]byte, int64, error)
}

// SessionRepository 会话仓储接口
type SessionRepository interface {
	// GetCurrentSession 获取当前会话
	GetCurrentSession(ctx context.Context) (*NpfsSession, error)

	// RefreshSession 刷新会话
	RefreshSession(ctx context.Context) error

	// CloseSession 关闭会话
	CloseSession(ctx context.Context) error
}
