package domain

import (
	"app/internal/shared/entity"
	"app/internal/shared/valueobject"
	"time"
)

// File NPFS 文件实体
type File struct {
	entity.BaseEntity // 继承基础实体字段

	Path     string                 `json:"path"`                        // 文件路径
	Name     string                 `json:"name"`                        // 文件名
	Size     int64                  `json:"size"`                        // 文件大小
	Hash     string                 `json:"hash"`                        // 文件哈希
	IsDir    bool                   `json:"is_dir"`                      // 是否为目录
	Metadata map[string]interface{} `json:"metadata,omitempty" gorm:"-"` // 元数据（不入库）
}

// TableName 指定表名
func (File) TableName() string {
	return "npfs_file"
}

// Directory NPFS 目录实体
type Directory struct {
	entity.BaseEntity

	Path     string      `json:"path"`           // 目录路径
	Name     string      `json:"name"`           // 目录名
	ParentID int64       `json:"parent_id"`      // 父目录 ID
	Files    []File      `json:"files" gorm:"-"` // 文件列表（非数据库字段）
	Dirs     []Directory `json:"dirs" gorm:"-"`  // 子目录列表（非数据库字段）
}

// TableName 指定表名
func (Directory) TableName() string {
	return "npfs_directory"
}

// NpfsSession NPFS 会话实体
type NpfsSession struct {
	entity.BaseEntity

	SessionID   string                 `json:"session_id"`            // 会话 ID
	ClientURL   string                 `json:"client_url"`            // 客户端 URL
	LoginTime   time.Time              `json:"login_time"`            // 登录时间
	ExpiredAt   time.Time              `json:"expired_at"`            // 过期时间
	Coordinates valueobject.Coordinate `json:"coordinates,omitempty"` // 坐标信息（可选）
}

// TableName 指定表名
func (NpfsSession) TableName() string {
	return "npfs_session"
}
