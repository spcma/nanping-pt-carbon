package http

import (
	"app/internal/module/carbonreport/application"
	http2 "app/internal/platform/http"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// FileHandler 文件 HTTP 处理器
type FileHandler struct {
	appService *application.FileApplicationService
}

// NewFileHandler 创建文件处理器
func NewFileHandler(appService *application.FileApplicationService) *FileHandler {
	return &FileHandler{
		appService: appService,
	}
}

// CheckDirHandler 检查目录 handler
func (h *FileHandler) CheckDirHandler(c *gin.Context) {
	var req application.CheckDirRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误：" + err.Error()})
		return
	}

	exists, err := h.appService.CheckDirExists(http2.Ctx(c), req.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "检查失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"path":   req.Path,
			"exists": exists,
		},
	})
}

// CreateDirHandler 创建目录 handler
func (h *FileHandler) CreateDirHandler(c *gin.Context) {
	var req application.CreateDirRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误：" + err.Error()})
		return
	}

	err := h.appService.CreateDirectory(http2.Ctx(c), req.Path, req.Recursive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"path":    req.Path,
			"created": true,
		},
	})
}

// ListDirHandler 列出目录 handler
func (h *FileHandler) ListDirHandler(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少参数：path"})
		return
	}

	files, err := h.appService.ListDirectory(http2.Ctx(c), path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "列出失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"path":  path,
			"files": files,
		},
	})
}

// DeleteDirHandler 删除目录 handler
func (h *FileHandler) DeleteDirHandler(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少参数：path"})
		return
	}

	recursive := c.Query("recursive") == "true"
	force := c.Query("force") == "true"

	err := h.appService.DeleteFile(http2.Ctx(c), path, recursive, force)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"path":    path,
			"deleted": true,
		},
	})
}

// ReadFileHandler 读取文件 handler
func (h *FileHandler) ReadFileHandler(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少参数：path"})
		return
	}

	data, size, err := h.appService.ReadFile(http2.Ctx(c), path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "读取失败：" + err.Error()})
		return
	}

	content := ""
	if size < 10000 {
		content = string(data)
	} else {
		content = "(文件太大，仅显示前 10000 字节)"
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"path":    path,
			"size":    size,
			"content": content,
		},
	})
}

// SaveFileHandler 保存文件 handler
func (h *FileHandler) SaveFileHandler(c *gin.Context) {
	var req application.SaveFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误：" + err.Error()})
		return
	}

	ipfsid, err := h.appService.SaveContent(http2.Ctx(c), req.Content, req.Dir, req.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"ipfsid": ipfsid,
			"path":   req.Dir + "/" + req.Filename,
		},
	})
}

// UploadFileHandler 上传文件 handler
func (h *FileHandler) UploadFileHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少文件参数"})
		return
	}

	dir := c.PostForm("dir")
	if dir == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少参数：dir"})
		return
	}

	filename := c.PostForm("filename")
	if filename == "" {
		filename = file.Filename
	}

	// 保存上传的文件到临时位置
	tmpPath := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveUploadedFile(file, tmpPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "保存临时文件失败：" + err.Error()})
		return
	}
	defer os.Remove(tmpPath)

	// 保存到 NPFS
	ipfsid, err := h.appService.SaveLocalFile(http2.Ctx(c), tmpPath, dir, filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "上传失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"ipfsid": ipfsid,
			"path":   dir + "/" + filename,
		},
	})
}

// DownloadFileHandler 下载文件 handler
func (h *FileHandler) DownloadFileHandler(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少参数：path"})
		return
	}

	filename := c.Query("filename")
	if filename == "" {
		filename = filepath.Base(path)
	}

	data, _, err := h.appService.ReadFile(http2.Ctx(c), path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "读取失败：" + err.Error()})
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", len(data)))

	c.Data(http.StatusOK, "application/octet-stream", data)
}

// DeleteFileHandler 删除文件 handler
func (h *FileHandler) DeleteFileHandler(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少参数：path"})
		return
	}

	force := c.Query("force") == "true"

	err := h.appService.DeleteFile(http2.Ctx(c), path, false, force)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败：" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"path":    path,
			"deleted": true,
		},
	})
}
