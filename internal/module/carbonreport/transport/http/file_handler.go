package http

import (
	"app/internal/module/carbonreport/application"
	http2 "app/internal/platform/http"
	"app/internal/platform/http/response"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// FileHandler 文件处理器
type FileHandler struct {
	appService *application.FileAppService
}

// NewFileHandler 创建文件处理器
func NewFileHandler(appService *application.FileAppService) *FileHandler {
	return &FileHandler{
		appService: appService,
	}
}

// CheckDir 检查目录
func (h *FileHandler) CheckDir(c *gin.Context) {
	var cmd application.CheckDirCommand
	if err := c.ShouldBind(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.appService.CheckDir(http2.Ctx(c), cmd)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// CreateDir 创建目录
func (h *FileHandler) CreateDir(c *gin.Context) {
	var cmd application.CreateDirCommand
	if err := c.ShouldBind(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.appService.CreateDir(http2.Ctx(c), cmd)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// ListDir 列出目录
func (h *FileHandler) ListDir(c *gin.Context) {
	var cmd application.ListDirCommand
	if err := c.ShouldBindQuery(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.appService.ListDir(http2.Ctx(c), cmd)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteFile 删除文件
func (h *FileHandler) DeleteFile(c *gin.Context) {
	var cmd application.DeleteFileCommand
	if err := c.ShouldBindQuery(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.appService.DeleteFile(http2.Ctx(c), cmd)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// ReadFile 读取文件
func (h *FileHandler) ReadFile(c *gin.Context) {
	var cmd application.ReadFileCommand
	if err := c.ShouldBindQuery(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.appService.ReadFile(http2.Ctx(c), cmd)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// SaveFile 保存文件
func (h *FileHandler) SaveFile(c *gin.Context) {
	var cmd application.SaveFileCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.appService.SaveFile(http2.Ctx(c), cmd)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// UploadFile 上传文件
func (h *FileHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "缺少文件参数")
		return
	}

	dir := c.PostForm("dir")
	if dir == "" {
		response.Error(c, http.StatusBadRequest, "缺少参数：dir")
		return
	}

	filename := c.PostForm("filename")

	// 保存上传文件到临时位置
	tmpPath := filepath.Join("/tmp", file.Filename)
	if err := c.SaveUploadedFile(file, tmpPath); err != nil {
		response.Error(c, http.StatusInternalServerError, "保存临时文件失败："+err.Error())
		return
	}

	cmd := application.UploadFileCommand{
		Dir:      dir,
		Filename: filename,
	}

	result, err := h.appService.UploadFile(http2.Ctx(c), tmpPath, cmd)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// DownloadFile 下载文件
func (h *FileHandler) DownloadFile(c *gin.Context) {
	var cmd application.DownloadFileCommand
	if err := c.ShouldBindQuery(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	data, filename, err := h.appService.DownloadFile(http2.Ctx(c), cmd)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename*=UTF-8''"+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", string(len(data)))

	c.Data(http.StatusOK, "application/octet-stream", data)
}
