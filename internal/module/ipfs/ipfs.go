package ipfs

import (
	"app/internal/rpc"
	"app/internal/shared/logger"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Service IPFS 服务
type Service struct {
	client  *rpc.LApiStub
	session string
}

// NewService 创建 IPFS 服务
func NewService(client *rpc.LApiStub, session string) *Service {
	return &Service{
		client:  client,
		session: session,
	}
}

// FileResponse 文件响应
type FileResponse struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Content string `json:"content,omitempty"`
}

// DirResponse 目录项
type DirResponse struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Type string `json:"type"` // file or dir
	Size int64  `json:"size"`
}

// CheckDirRequest 检查目录请求
type CheckDirRequest struct {
	Path string `json:"path" form:"path"`
}

// CreateDirRequest 创建目录请求
type CreateDirRequest struct {
	Path string `json:"path" form:"path"`
}

// ListDirRequest 列出目录请求
type ListDirRequest struct {
	Path string `form:"path"`
}

// DeleteFileRequest 删除文件请求
type DeleteFileRequest struct {
	Path string `form:"path"`
}

// ReadFileRequest 读取文件请求
type ReadFileRequest struct {
	Path string `form:"path"`
}

// SaveFileRequest 保存文件请求
type SaveFileRequest struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// UploadFileCommand 上传文件命令
type UploadFileCommand struct {
	Dir      string
	Filename string
}

// DownloadFileRequest 下载文件请求
type DownloadFileRequest struct {
	Path string `form:"path"`
}

// CheckDir 检查目录
func (s *Service) CheckDir(c *gin.Context) {
	var req CheckDirRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 实现 IPFS 目录检查逻辑
	logger.Info("ipfs", "check dir",
		zap.String("path", req.Path),
	)

	c.JSON(http.StatusOK, gin.H{
		"exists": true,
		"path":   req.Path,
	})
}

// CreateDir 创建目录
func (s *Service) CreateDir(c *gin.Context) {
	var req CreateDirRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 实现 IPFS 创建目录逻辑
	logger.Info("ipfs", "create dir",
		zap.String("path", req.Path),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"path":    req.Path,
	})
}

// ListDir 列出目录
func (s *Service) ListDir(c *gin.Context) {
	var req ListDirRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 实现 IPFS 列出目录逻辑
	logger.Info("ipfs", "list dir",
		zap.String("path", req.Path),
	)

	c.JSON(http.StatusOK, gin.H{
		"list": []DirResponse{},
	})
}

// DeleteFile 删除文件
func (s *Service) DeleteFile(c *gin.Context) {
	var req DeleteFileRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 实现 IPFS 删除文件逻辑
	logger.Info("ipfs", "delete file",
		zap.String("path", req.Path),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// ReadFile 读取文件
func (s *Service) ReadFile(c *gin.Context) {
	var req ReadFileRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 实现 IPFS 读取文件逻辑
	logger.Info("ipfs", "read file",
		zap.String("path", req.Path),
	)

	c.JSON(http.StatusOK, gin.H{
		"content": "",
		"path":    req.Path,
	})
}

// SaveFile 保存文件
func (s *Service) SaveFile(c *gin.Context) {
	var req SaveFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 实现 IPFS 保存文件逻辑
	logger.Info("ipfs", "save file",
		zap.String("path", req.Path),
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"path":    req.Path,
	})
}

// UploadFile 上传文件
func (s *Service) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少文件参数"})
		return
	}

	dir := c.PostForm("dir")
	if dir == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少参数：dir"})
		return
	}

	filename := c.PostForm("filename")
	if filename == "" {
		filename = file.Filename
	}

	// 保存上传文件到临时位置
	tmpPath := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveUploadedFile(file, tmpPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存临时文件失败：" + err.Error()})
		return
	}
	defer os.Remove(tmpPath)

	// 读取文件内容
	fileData, err := os.ReadFile(tmpPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文件失败：" + err.Error()})
		return
	}

	// TODO: 实现 IPFS 上传文件逻辑
	logger.Info("ipfs", "upload file",
		zap.String("dir", dir),
		zap.String("filename", filename),
		zap.Int("size", len(fileData)),
	)

	fullPath := strings.TrimSuffix(dir, "/") + "/" + filename

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"path":    fullPath,
		"size":    len(fileData),
	})
}

// DownloadFile 下载文件
func (s *Service) DownloadFile(c *gin.Context) {
	var req DownloadFileRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 实现 IPFS 下载文件逻辑
	logger.Info("ipfs", "download file",
		zap.String("path", req.Path),
	)

	// 示例：返回空文件
	filename := filepath.Base(req.Path)
	c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.Data(http.StatusOK, "application/octet-stream", []byte{})
}

// CreateFsClient 创建文件系统客户端
func CreateFsClient() (*rpc.LApiStub, string, error) {
	strPPT, err := rpc.GetLocalPassport(4080, 24)
	if err != nil {
		return nil, "", err
	}

	client := rpc.InitLApiStubByUrl("127.0.0.1:4080")

	loginReply, err := client.LoginWithPPT(strPPT)
	if err != nil {
		return nil, "", err
	}

	return client, loginReply.Sid, nil
}
