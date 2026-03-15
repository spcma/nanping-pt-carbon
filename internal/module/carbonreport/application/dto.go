package application

import (
	"app/internal/module/carbonreport/domain"
)

// FileDTO 文件数据传输对象
type FileDTO struct {
	ID        int64  `json:"id"`
	Path      string `json:"path"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	Hash      string `json:"hash"`
	IsDir     bool   `json:"is_dir"`
	IsDeleted bool   `json:"is_deleted,omitempty"`
}

// DirectoryDTO 目录数据传输对象
type DirectoryDTO struct {
	ID    int64     `json:"id"`
	Path  string    `json:"path"`
	Name  string    `json:"name"`
	Files []FileDTO `json:"files"`
}

// CheckDirRequest 检查目录请求
type CheckDirRequest struct {
	Path string `json:"path" form:"path" binding:"required"`
}

// CheckDirResponse 检查目录响应
type CheckDirResponse struct {
	Path   string `json:"path"`
	Exists bool   `json:"exists"`
}

// CreateDirRequest 创建目录请求
type CreateDirRequest struct {
	Path      string `json:"path" form:"path" binding:"required"`
	Recursive bool   `json:"recursive" form:"recursive"`
}

// CreateDirResponse 创建目录响应
type CreateDirResponse struct {
	Path    string `json:"path"`
	Created bool   `json:"created"`
}

// ListDirRequest 列出目录请求
type ListDirRequest struct {
	Path string `json:"path" form:"path" query:"path" binding:"required"`
}

// ListDirResponse 列出目录响应
type ListDirResponse struct {
	Path  string    `json:"path"`
	Files []FileDTO `json:"files"`
}

// DeleteFileRequest 删除文件请求
type DeleteFileRequest struct {
	Path      string `json:"path" form:"path" query:"path" binding:"required"`
	Recursive bool   `json:"recursive" form:"recursive"`
	Force     bool   `json:"force" form:"force"`
}

// DeleteFileResponse 删除文件响应
type DeleteFileResponse struct {
	Path    string `json:"path"`
	Deleted bool   `json:"deleted"`
}

// ReadFileRequest 读取文件请求
type ReadFileRequest struct {
	Path string `json:"path" form:"path" query:"path" binding:"required"`
}

// ReadFileResponse 读取文件响应
type ReadFileResponse struct {
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	Content string `json:"content,omitempty"`
	Data    []byte `json:"data,omitempty"`
}

// SaveFileRequest 保存文件请求
type SaveFileRequest struct {
	Content  string `json:"content" binding:"required"`
	Dir      string `json:"dir" binding:"required"`
	Filename string `json:"filename" binding:"required"`
}

// SaveFileResponse 保存文件响应
type SaveFileResponse struct {
	IPFSID string `json:"ipfsid"`
	Path   string `json:"path"`
}

// UploadFileRequest 上传文件请求
type UploadFileRequest struct {
	Dir      string `form:"dir" binding:"required"`
	Filename string `form:"filename"`
}

// UploadFileResponse 上传文件响应
type UploadFileResponse struct {
	IPFSID string `json:"ipfsid"`
	Path   string `json:"path"`
}

// DownloadFileRequest 下载文件请求
type DownloadFileRequest struct {
	Path     string `json:"path" form:"path" query:"path" binding:"required"`
	Filename string `json:"filename" form:"filename" query:"filename"`
}

// ConvertFileToDTO 将领域文件转换为 DTO
func ConvertFileToDTO(file *domain.File) FileDTO {
	return FileDTO{
		ID:        file.Id, // 使用 BaseEntity 的 Id
		Path:      file.Path,
		Name:      file.Name,
		Size:      file.Size,
		Hash:      file.Hash,
		IsDir:     file.IsDir,
		IsDeleted: file.IsDeleted(),
	}
}

// ConvertFilesToDTOs 批量转换
func ConvertFilesToDTOs(files []domain.File) []FileDTO {
	result := make([]FileDTO, 0, len(files))
	for _, file := range files {
		result = append(result, ConvertFileToDTO(&file))
	}
	return result
}

// ConvertDirectoryToDTO 将目录转换为 DTO
func ConvertDirectoryToDTO(dir *domain.Directory) DirectoryDTO {
	return DirectoryDTO{
		ID:    dir.Id,
		Path:  dir.Path,
		Name:  dir.Name,
		Files: ConvertFilesToDTOs(dir.Files),
	}
}
