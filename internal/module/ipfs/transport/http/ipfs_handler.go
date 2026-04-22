package http

import (
	"app/internal/module/ipfs/application"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dromara/carbon/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IpfsHandler struct {
	appService *application.Service
}

func NewIpfsHandler(appService *application.Service) *IpfsHandler {
	return &IpfsHandler{
		appService: appService,
	}
}

func (h *IpfsHandler) ListDir(c *gin.Context) {
	dir := c.Query("dir")
	if dir == "" {
		response.BadRequest(c, "dir is required")
		return
	}

	listDir, err := h.appService.ListDir(platform_http.Ctx(c), dir)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, listDir)
}

func (h *IpfsHandler) CreateDir(c *gin.Context) {

	type CreateDirDto struct {
		Dir string `json:"dir"`
	}

	var dto CreateDirDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if dto.Dir == "" {
		response.BadRequest(c, "dir is required")
		return
	}

	err := h.appService.CreateDir(dto.Dir)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, "create dir success")
}

func (h *IpfsHandler) DeleteFile(c *gin.Context) {

	type DeleteFileDto struct {
		Path string `json:"path"`
	}

	var dto DeleteFileDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if dto.Path == "" {
		response.BadRequest(c, "path is required")
		return
	}

	err := h.appService.DeleteFile(platform_http.Ctx(c), dto.Path)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, "delete file success")
}

func (h *IpfsHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		return
	}

	dir := c.PostForm("dir")
	if dir == "" {
		return
	}

	filename := c.PostForm("filename")
	if filename == "" {
		filename = file.Filename
	}

	// 保存上传的文件到临时位置
	tmpPath := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveUploadedFile(file, tmpPath); err != nil {
		logger.IpfsL.Error("save uploaded file failed", zap.String("tmpPath", tmpPath), zap.Error(err))
		return
	}
	defer os.Remove(tmpPath)

	ipfsid, err := h.appService.UploadFile(tmpPath, dir, filename)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, ipfsid)
}

func (h *IpfsHandler) DownloadFile(c *gin.Context) {

	type downloadDto struct {
		Path     string `json:"path" form:"path"`
		Filename string `json:"filename" form:"filename"`
	}

	var dto downloadDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if dto.Path == "" {
		response.BadRequest(c, "path is required")
		return
	}

	if dto.Filename == "" {
		dto.Filename = filepath.Base(dto.Path)
	}

	data, _, err := h.appService.ReadFileFromIpfs(dto.Path)
	if err != nil {
		response.InternalError(c, "下载失败："+err.Error())
		return
	}

	logger.IpfsL.Info("download file",
		zap.String("path", dto.Path),
		zap.String("filename", dto.Filename),
		zap.Any("data", string(data)))

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", dto.Filename))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Length", fmt.Sprintf("%d", len(data)))

	c.Data(http.StatusOK, "application/octet-stream", data)
}

// Stat 获取目录/文件信息
func (h *IpfsHandler) Stat(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		response.BadRequest(c, "请指定目录")
		return
	}

	fileStat, err := h.appService.FileStat(path)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, fileStat)
}

// Read 读取ipfs文件
func (h *IpfsHandler) Read(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		response.BadRequest(c, "请指定目录")
		return
	}

	bytes, count, err := h.appService.ReadFileFromIpfs(path)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	logger.IpfsL.Info("readIpfs", zap.Any("count", count), zap.Any("bytes", bytes))

	response.Success(c, string(bytes))
}

// ScanDir 递归扫描目录
func (h *IpfsHandler) ScanDir(c *gin.Context) {
	type ScanDirDto struct {
		RootDir string `json:"rootDir" form:"rootDir"`
	}

	var dto ScanDirDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if dto.RootDir == "" {
		response.BadRequest(c, "请指定根目录")
		return
	}

	go func() {
		ctx := context.Background()
		result, err := h.appService.ScanDir(ctx, dto.RootDir)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}
		logger.IpfsL.Info("scanDir completed", zap.String("rootDir", dto.RootDir), zap.Any("result", result))
	}()

	response.Success(c, "scanDir task running...")
}

func (h *IpfsHandler) CalcDir(c *gin.Context) {
	type CalcDirDto struct {
		RootDir string `json:"rootDir" form:"rootDir"` // 要扫描的根目录，如 "/aibk/26/03/27"
		Date    string `json:"date" form:"date"`       // 日期，格式 "2026-03-27"
	}

	var dto CalcDirDto
	if err := c.ShouldBindQuery(&dto); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if dto.RootDir == "" {
		response.BadRequest(c, "请指定目录")
		return
	}

	if dto.Date == "" {
		response.BadRequest(c, "请指定日期")
		return
	}

	go func() {
		ctx := context.Background()
		turnover, err := h.appService.CalcDir(ctx, dto.RootDir, dto.Date)
		if err != nil {
			logger.IpfsL.Error("calcDir failed", zap.String("rootDir", dto.RootDir), zap.String("date", dto.Date), zap.Error(err))
			return
		}

		logger.IpfsL.Info("calcDir completed", zap.String("rootDir", dto.RootDir), zap.String("date", dto.Date), zap.Float64("turnover", turnover))
	}()

	response.Success(c, "计算任务已启动，请稍后查看结果")
}

func (h *IpfsHandler) SaveContentTest(c *gin.Context) {

	//value := c.Query("force")
	//force := cast.ToBool(value)

	now := carbon.Now().StartOfDay()

	//	 保存周转量结果到文件中 /tmpp/26/03/14
	saveDir := fmt.Sprintf("%s/%s/%s/%s", "/tmpp", "27", "03", "14")

	filename := fmt.Sprintf("%s.txt", now.Format(carbon.ShortDateFormat))

	//if force {
	//	h.appService.Remove()
	//}

	ipfsid, err := h.appService.SaveContentToIpfs("hello world", saveDir, filename)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, ipfsid)
}
