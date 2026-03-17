package ipfs

import (
	"app/internal/module/ipfs/rpc"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// Service IPFS 服务
type Service struct {
	db      *gorm.DB
	client  *rpc.LApiStub
	session string
}

// NewService 创建 IPFS 服务
func NewService(db *gorm.DB, client *rpc.LApiStub, session string) *Service {
	return &Service{
		db:      db,
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

// UploadFileCommand 上传文件命令
type UploadFileCommand struct {
	Dir      string
	Filename string
}

// CheckDirRequest 检查目录请求
type CheckDirRequest struct {
	Path string `json:"path" form:"path"`
}

// CheckDir 检查目录
func (s *Service) CheckDir(c *gin.Context) {
	var req CheckDirRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.IpfsLogger.Info("check dir", zap.String("path", req.Path))

	c.JSON(http.StatusOK, gin.H{
		"exists": true,
		"path":   req.Path,
	})
}

// CreateDirRequest 创建目录请求
type CreateDirRequest struct {
	Path string `json:"path" form:"path"`
}

// CreateDir 创建目录
func (s *Service) CreateDir(c *gin.Context) {
	var req CreateDirRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.IpfsLogger.Info("create dir", zap.String("path", req.Path))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"path":    req.Path,
	})
}

// ListDirRequest 列出目录请求
type ListDirRequest struct {
	Path string `form:"path"`
}

// DirResponse 目录项
type DirResponse struct {
	Name      string
	Hash      string
	Size      uint64
	Type      int    // 1 dir 0 file
	FileType  string // .jpg .txt
	Timestamp int64  // 时间戳
	Seq       int    // 序号
}

// ListDir 列出目录
func (s *Service) ListDir(c *gin.Context) {
	var req ListDirRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	logger.IpfsLogger.Info("list dir", zap.String("path", req.Path))

	lsLinks, err := s.client.FilesLs(s.session, req.Path)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var list []DirResponse
	for _, link := range lsLinks {

		dir := DirResponse{
			Name: link.Name,
			Hash: link.Hash,
			Size: link.Size,
			Type: link.Type,
		}

		ext := filepath.Ext(link.Name)

		newStr := strings.ReplaceAll(link.Name, "image_", "")
		newStr = strings.ReplaceAll(link.Name, "gps_", "")
		newStr = strings.ReplaceAll(newStr, ext, "")

		if ext == ".jpg" {
			split := strings.Split(newStr, "_")
			if len(split) >= 2 {
				dir.Timestamp = cast.ToInt64(split[0])
				dir.Seq = cast.ToInt(split[1])
			} else {
				logger.IpfsLogger.Error("parse time error", zap.String("file", link.Name))
			}
		} else {
			dir.Timestamp = cast.ToInt64(newStr)
		}

		list = append(list, dir)
	}

	response.Success(c, list)
}

func (s *Service) HandleWithDir(c *gin.Context) {
	type Param struct {
		Dir        string `json:"dir" form:"dir"` // 目录
		Year       string `json:"year" form:"year"`
		Month      string `json:"month" form:"month"`
		Day        string `json:"day" form:"day"`
		DeviceCode string `json:"device_code" form:"device_code"`
	}

	var req Param
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	//	开始扫描目录
	fullDir := fmt.Sprintf("%s/%s/%s/%s/%s", req.Dir, req.Year, req.Month, req.Day, req.DeviceCode)

	logger.IpfsLogger.Info("handle with dir", zap.String("dir", req.Dir),
		zap.String("year", req.Year),
		zap.String("month", req.Month),
		zap.String("day", req.Day),
		zap.String("device_code", req.DeviceCode),
		zap.String("full_dir", fullDir),
	)

	lsLinks, err := s.client.FilesLs(s.session, fullDir)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	var list []DirResponse
	for _, link := range lsLinks {

		if link.Type == 1 {
			logger.IpfsLogger.Info("skip dir", zap.String("dir", link.Name))
			continue
		}

		fullPath := fmt.Sprintf("%s/%s", fullDir, link.Name)

		localPath := "./tempfile/" + link.Name
		err := s.SaveFileToLocal(fullPath, localPath)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		rec, err := parseFile(localPath)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		logger.IpfsLogger.Info("parse file", zap.String("file", link.Name),
			zap.Int("count", len(rec)),
		)

		// 计算里程并保存到文件
		calculator := NewDistanceCalculator()
		totalDistance, err := calculator.CalculateTotalDistance(rec)
		if err == nil && len(rec) > 0 {
			summary := calculator.CalculateSummary(rec)
			logger.IpfsLogger.Info("distance calculation",
				zap.String("file", link.Name),
				zap.Float64("total_distance_m", totalDistance),
				zap.Float64("total_distance_km", summary.TotalDistanceKm),
				zap.Int("point_count", summary.PointCount),
				zap.Float64("avg_speed_kmh", summary.AverageSpeed),
			)

			// 将里程结果写入新文件
			err = s.saveDistanceResult(link.Name, summary)
			if err != nil {
				logger.IpfsLogger.Warn("save distance result failed",
					zap.String("file", link.Name),
					zap.Error(err),
				)
			}
		}

		for _, record := range rec {
			logger.IpfsLogger.Info("parse file", zap.String("file", link.Name),
				zap.Time("timestamp", record.Timestamp),
				zap.Float64("lat", record.Lat),
				zap.Float64("lon", record.Lon),
				zap.Float64("value", record.Value),
			)
		}

		dir := DirResponse{
			Name: link.Name,
			Hash: link.Hash,
			Size: link.Size,
			Type: link.Type,
		}

		ext := filepath.Ext(link.Name)

		newStr := strings.ReplaceAll(link.Name, "image_", "")
		newStr = strings.ReplaceAll(link.Name, "gps_", "")
		newStr = strings.ReplaceAll(newStr, ext, "")

		if ext == ".jpg" {
			split := strings.Split(newStr, "_")
			if len(split) >= 2 {
				dir.Timestamp = cast.ToInt64(split[0])
				dir.Seq = cast.ToInt(split[1])
			} else {
				logger.IpfsLogger.Error("parse time error", zap.String("file", link.Name))
			}
		} else {
			dir.Timestamp = cast.ToInt64(newStr)
		}

		list = append(list, dir)
	}
}

// DeleteFile 删除文件
func (s *Service) DeleteFile(c *gin.Context) {
	var req DeleteFileRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.IpfsLogger.Info("delete file", zap.String("path", req.Path))

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

	logger.IpfsLogger.Info("read file", zap.String("path", req.Path))

	ext := filepath.Ext(req.Path)

	err := s.SaveFileToLocal(req.Path, "./tmp"+ext)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	rec, err := parseFile("./tmp" + ext)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 计算里程
	calculator := NewDistanceCalculator()
	_, _ = calculator.CalculateTotalDistance(rec)
	summary := calculator.CalculateSummary(rec)

	for _, record := range rec {
		logger.IpfsLogger.Info("read file", zap.Any("record", record))
	}

	c.JSON(http.StatusOK, gin.H{
		"content":          rec,
		"path":             req.Path,
		"distance_summary": summary,
	})
}

// ==================== 文件读取相关 ====================

// ReadFileFromNpfs 从 NPFS 读取文件数据
// filePath: NPFS 文件路径（如：/np_storage/1.jpg）
// data: 文件数据，size: 文件大小，err: 错误信息
func (s *Service) ReadFileFromNpfs(filePath string) ([]byte, int64, error) {
	// 打开文件 URL
	fsid, err := s.client.MMOpenUrl(s.session, filePath)
	if err != nil {
		return nil, 0, err
	}
	defer s.client.MMClose(fsid)

	// 获取文件大小
	size, err := s.client.MFGetSize(fsid)
	if err != nil {
		return nil, 0, err
	}

	// 读取文件数据
	data, err := s.client.MFGetData(fsid, 0, int(size))
	if err != nil {
		return nil, 0, err
	}

	return data, size, nil
}

// SaveFileToLocal 将 NPFS 文件保存到本地
// filePath: NPFS 文件路径
// localPath: 本地保存路径
func (s *Service) SaveFileToLocal(filePath, localPath string) error {
	data, _, err := s.ReadFileFromNpfs(filePath)
	if err != nil {
		return err
	}

	err = os.WriteFile(localPath, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// ReadFile 读取文件
func (s *Service) readFile(filePath string) ([]byte, int64, error) {
	fsid, err := s.client.MMOpenUrl(s.session, filePath)
	if err != nil {
		return nil, 0, err
	}
	defer s.client.MMClose(fsid)

	size, err := s.client.MFGetSize(fsid)
	if err != nil {
		return nil, 0, err
	}

	data, err := s.client.MFGetData(fsid, 0, int(size))
	if err != nil {
		return nil, 0, err
	}

	return data, size, nil
}

// SaveFileRequest 保存文件请求
type SaveFileRequest struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// SaveFile 保存文件
func (s *Service) SaveFile(c *gin.Context) {
	var req SaveFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.IpfsLogger.Info("save file", zap.String("path", req.Path))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"path":    req.Path,
	})
}

// DeleteFileRequest 删除文件请求
type DeleteFileRequest struct {
	Path string `form:"path"`
}

// ReadFileRequest 读取文件请求
type ReadFileRequest struct {
	Path string `form:"path"`
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

	logger.IpfsLogger.Info("upload file", zap.String("dir", dir), zap.String("filename", filename), zap.Int("size", len(fileData)))

	fullPath := strings.TrimSuffix(dir, "/") + "/" + filename

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"path":    fullPath,
		"size":    len(fileData),
	})
}

// DownloadFileRequest 下载文件请求
type DownloadFileRequest struct {
	Path string `form:"path"`
}

// DownloadFile 下载文件
func (s *Service) DownloadFile(c *gin.Context) {
	var req DownloadFileRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.IpfsLogger.Info("download file", zap.String("path", req.Path))

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

// saveDistanceResult 将里程计算结果保存到文件
// fileName: 原文件名（如：gps_20260314142635.txt）
// summary: 里程汇总信息
func (s *Service) saveDistanceResult(fileName string, summary DistanceSummary) error {
	// 生成结果文件路径，例如：./distance_result/gps_20260314142635_distance.txt
	resultDir := "./distance_result"

	// 确保目录存在
	if err := os.MkdirAll(resultDir, os.ModePerm); err != nil {
		return err
	}

	// 生成结果文件名
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	resultFileName := baseName + "_distance.txt"
	resultPath := filepath.Join(resultDir, resultFileName)

	// 格式化输出内容
	content := fmt.Sprintf("%.6f\n", summary.TotalDistanceKm)

	// 写入文件
	err := os.WriteFile(resultPath, []byte(content), 0644)
	if err != nil {
		return err
	}

	logger.IpfsLogger.Info("distance result saved",
		zap.String("source_file", fileName),
		zap.String("result_file", resultPath),
		zap.Float64("distance_km", summary.TotalDistanceKm),
	)

	return nil
}
