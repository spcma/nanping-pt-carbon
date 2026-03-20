package application

import (
	carbonreportday_application "app/internal/module/carbonreportday/application"
	carbonreportday_infrastructure "app/internal/module/carbonreportday/infrastructure"
	"app/internal/module/ipfs/infrastructure"
	"app/internal/module/ipfs/rpc"
	"app/internal/platform/http/response"
	"app/internal/shared/entity"
	"app/internal/shared/logger"
	"app/internal/shared/timeutil"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dromara/carbon/v2"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

// Service IPFS 服务
type Service struct {
	db                        *gorm.DB
	remoteDB                  *gorm.DB
	client                    *rpc.LApiStub
	session                   string
	ipfsDetailAppService      *IpfsDetailAppService
	carbonReportDayAppService *carbonreportday_application.CarbonReportDayAppService
}

// NewService 创建 IPFS 服务
func NewService(db *gorm.DB, remoteDB *gorm.DB, client *rpc.LApiStub, session string) *Service {
	// 初始化仓储和應用服務
	ipfsDetailRepo := infrastructure.NewIpfsDetailRepository(db)
	ipfsDetailAppService := NewIpfsDetailAppService(ipfsDetailRepo)

	// 初始化碳报告应用服务
	carbonReportDayRepo := carbonreportday_infrastructure.NewCarbonReportDayRepository(db)
	carbonReportDayAppService := carbonreportday_application.NewCarbonReportDayAppService(carbonReportDayRepo)

	return &Service{
		db:                        db,
		remoteDB:                  remoteDB,
		client:                    client,
		session:                   session,
		ipfsDetailAppService:      ipfsDetailAppService,
		carbonReportDayAppService: carbonReportDayAppService,
	}
}

// DirDto 目录项
type DirDto struct {
	Name      string
	Hash      string
	Size      uint64
	Type      int    // 1 dir 0 file
	FileType  string // .jpg .txt
	Timestamp int64  // 时间戳
	Seq       int    // 序号
}

func (s *Service) CheckDir(ctx context.Context, dir string) bool {
	stat, err := s.client.FilesStat(s.session, dir)
	if err != nil {
		if fileNotExist(err) {
			logger.IpfsLogger.Warn("目录不存在", zap.String("dir", dir))
			return false
		}
		logger.IpfsLogger.Error("发生了错误", zap.Error(err))
		return false
	}

	logger.IpfsLogger.Info("检查目录",
		zap.String("dir", dir),
		zap.Any("stat", stat),
		zap.Any("hash", stat.Hash),
		zap.Any("size", stat.Size),
		zap.Any("cumulativeSize", stat.CumulativeSize),
		zap.Any("blocks", stat.Blocks),
		zap.Any("type", stat.Type),
		zap.Any("withLocality", stat.WithLocality),
		zap.Any("local", stat.Local),
		zap.Any("sizeLocal", stat.SizeLocal),
	)

	return true
}

// CreateDir 创建目录
func (s *Service) CreateDir(ctx context.Context, dir string) error {
	if dir == "" {
		return errors.New("目录为空")
	}

	err := s.client.FilesMkdir(s.session, dir, true)
	if err != nil {
		logger.IpfsLogger.Error("创建目录失败", zap.String("dir", dir), zap.Error(err))
		return err
	}

	return nil
}

// ListDir 列出目录
func (s *Service) ListDir(ctx context.Context, dir string) ([]*DirDto, error) {
	if dir == "" {
		return nil, errors.New("目录为空")
	}

	lsLinks, err := s.client.FilesLs(s.session, dir)
	if err != nil {
		logger.IpfsLogger.Error("列出目录失败", zap.String("dir", dir), zap.Error(err))
		return nil, err
	}

	var list []*DirDto
	for _, link := range lsLinks {

		dirDto := &DirDto{
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
				dirDto.Timestamp = cast.ToInt64(split[0])
				dirDto.Seq = cast.ToInt(split[1])
			} else {
				logger.IpfsLogger.Error("parse time error", zap.String("file", link.Name))
			}
		} else {
			dirDto.Timestamp = cast.ToInt64(newStr)
		}

		list = append(list, dirDto)
	}

	return list, nil
}

// DeleteFileRequest 删除文件请求
type DeleteFileRequest struct {
	Path      string `form:"path"`
	Recursive bool   `form:"recursive"` // 递归
	Force     bool   `form:"force"`     // 强制删除
}

// DeleteFile 删除文件
func (s *Service) DeleteFile(ctx context.Context, path string, recursive, force bool) error {
	err := s.client.FilesRm(s.session, path, recursive, force)
	if err != nil {
		logger.IpfsLogger.Error("delete file error", zap.Error(err))
		return err
	}

	return nil
}

// SaveContent 保存内容到文件
func (s *Service) SaveContent(ctx context.Context, content, fsDir, filename string) (string, error) {
	// 打开临时文件
	fsid, err := s.client.MFOpenTempFile(s.session)
	if err != nil {
		return "", err
	}

	// 写入数据
	_, err = s.client.MFSetData(fsid, []byte(content), 0)
	if err != nil {
		return "", err
	}

	// 确保目录存在
	ok, err := s.MustDirExists(ctx, fsDir, true)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", errors.New("dir not exists & create failed")
	}

	// 保存到 NPFS
	nodePath := fsDir + "/" + filename
	ipfsid, err := s.client.MFTemp2Files(fsid, nodePath)
	if err != nil {
		return "", err
	}

	return ipfsid, nil
}

// ReadFile 读取文件
func (s *Service) ReadFile(ctx context.Context, path string) ([]byte, error) {

	logger.IpfsLogger.Info("read file", zap.String("path", path))

	ext := filepath.Ext(path)

	err := s.SaveFileToLocal(path, "./tmp"+ext)
	if err != nil {
		return nil, err
	}

	rec, err := parseFile("./tmp" + ext)
	if err != nil {
		return nil, err
	}

	// 计算里程
	calculator := NewDistanceCalculator()
	summary := calculator.CalculateSummary(rec)

	for _, record := range rec {
		logger.IpfsLogger.Info("read file", zap.Any("record", record))
	}

	logger.IpfsLogger.Info("readFile", zap.Any("summary", summary))

	return nil, nil
}

func (s *Service) ReadIpfs(ctx context.Context, path string) ([]byte, int64, error) {
	return s.readFileHandle(path)
}

// ReadFile 读取 ipfs 文件
func (s *Service) readFileHandle(filePath string) ([]byte, int64, error) {
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

// SaveLocalFile 保存本地文件到 NPFS
func (s *Service) SaveLocalFile(localPath, fsDir, filename string) (string, error) {
	// 打开临时文件
	fsid, err := s.client.MFOpenTempFile(s.session)
	if err != nil {
		return "", err
	}

	// 读取本地文件
	data, err := os.ReadFile(localPath)
	if err != nil {
		return "", err
	}

	// 写入数据
	_, err = s.client.MFSetData(fsid, data, 0)
	if err != nil {
		return "", err
	}

	// 确保目录存在
	//err = s.MustDirExists(fsDir, true)
	//if err != nil {
	//	return "", err
	//}

	// 保存到 NPFS
	nodePath := fsDir + "/" + filename
	ipfsid, err := s.client.MFTemp2Files(fsid, nodePath)
	if err != nil {
		return "", err
	}

	return ipfsid, nil
}

// MustDirExists 确保目录存在
func (s *Service) MustDirExists(ctx context.Context, path string, recursive bool) (bool, error) {
	if !s.CheckDir(ctx, path) {
		err := s.CreateDir(ctx, path)
		if err != nil {
			return false, err
		}

		return false, nil
	}

	return true, nil
}

func fileNotExist(err error) bool {
	if strings.Contains(err.Error(), "no link named") {
		return true
	}
	return false
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

	_ = data

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

// ReadFileRequest 读取文件请求
type ReadFileRequest struct {
	Path string `form:"path"`
}

// UploadFileCommand 上传文件命令
type UploadFileCommand struct {
	Dir      string
	Filename string
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

// saveIpfsDetailToDB 保存 IPFS 详情到数据库
// dir: 目录路径（用于提取设备编码）
// fileName: 文件名
// records: GPS 记录列表
// summary: 里程汇总信息
func (s *Service) saveIpfsDetailToDB(deviceCode, fileName string, timestamp int64, turnover float64, passenger int64, records []Record, summary DistanceSummary) error {
	if len(records) == 0 {
		return nil
	}

	collectionTime := timeutil.Now(carbon.Parse(cast.ToString(timestamp), carbon.Shanghai).StdTime())

	// 创建命令
	cmd := CreateIpfsDetailCommand{
		DeviceCode:     deviceCode,
		Filename:       fileName,
		CollectionTime: collectionTime.Format("2006-01-02 15:04:05"),
		TotalDistance:  summary.TotalDistanceKm,
		PointCount:     int64(summary.PointCount),
		UserID:         0, // 系统自动创建
	}

	// 检查文件是否已存在
	existingDetail, _ := s.ipfsDetailAppService.GetIpfsDetailByFilename(context.Background(), fileName)
	if existingDetail != nil {
		logger.IpfsLogger.Warn("ipfs detail already exists, skip saving",
			zap.String("file", fileName),
			zap.Int64("existing_id", existingDetail.Id),
		)
		return nil
	}

	// 保存到数据库
	_, err := s.ipfsDetailAppService.CreateIpfsDetail(context.Background(), cmd)
	return err
}

func (s *Service) H1(c *gin.Context) {
	var lll []*BusImageDetailCv
	err := s.remoteDB.WithContext(context.Background()).Table("bus_image_detail_cv").Find(&lll).Error
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, lll)
}

type BusImageDetailCv struct {
	entity.BaseEntity `map:"dive"`
	BusImageId        int64         `json:"busImageId" form:"busImageId" gorm:"column:bus_image_id" label:""`
	CollectionTime    timeutil.Time `json:"collectionTime" form:"collectionTime" gorm:"column:collection_time;default:now()" label:""`
	MergeSrcPath      string        `json:"mergeSrcPath" form:"mergeSrcPath" gorm:"column:merge_src_path" label:""`
	BaiduPath         string        `json:"baiduPath" form:"baiduPath" gorm:"column:baidu_path" label:""`
	BaiduResult       int64         `json:"baiduResult" form:"baiduResult" gorm:"column:baidu_result" label:""`
	CvType            string        `json:"cvType" form:"cvType" gorm:"column:cv_type" label:"识别类型 10 原图 20 预处理后图片"`
	DeviceCode        string        `json:"deviceCode" form:"deviceCode" gorm:"column:device_code" label:"设备编号"`
}

// Close 关闭连接
func (s *Service) Close() {
	if s.client != nil {
		s.client.Logout(s.session, "")
	}
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

	var list []DirDto
	for _, link := range lsLinks {

		if link.Type == 1 {
			logger.IpfsLogger.Info("skip dir", zap.String("dir", link.Name))
			continue
		}

		if strings.Contains(link.Name, ".jpg") {
			logger.IpfsLogger.Info("skip image", zap.String("file", link.Name))
			continue
		}

		fullPath := fmt.Sprintf("%s/%s", fullDir, link.Name)

		localPath := "./tempfile/" + link.Name
		err := s.SaveFileToLocal(fullPath, localPath)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		records, err := parseFile(localPath)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}

		logger.IpfsLogger.Info("parse file", zap.String("file", link.Name),
			zap.Int("count", len(records)),
		)

		// 计算里程并保存到数据库
		calculator := NewDistanceCalculator()
		summary := calculator.CalculateSummary(records)

		logger.IpfsLogger.Info("distance calculation",
			zap.String("file", link.Name),
			zap.Float64("total_distance_m", summary.TotalDistance),
			zap.Float64("total_distance_km", summary.TotalDistanceKm),
			zap.Int("point_count", summary.PointCount),
			zap.Float64("avg_speed_kmh", summary.AverageSpeed),
		)

		dir := DirDto{
			Name: link.Name,
			Hash: link.Hash,
			Size: link.Size,
			Type: link.Type,
		}

		ext := filepath.Ext(link.Name)

		newStr := strings.ReplaceAll(link.Name, "image_", "")
		newStr = strings.ReplaceAll(newStr, "gps_", "")
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

		// 保存里程计算结果到数据库
		//err = s.saveIpfsDetailToDB(fullDir, link.Name, dir.Timestamp, records, summary)
		//if err != nil {
		//	logger.IpfsLogger.Warn("save ipfs detail to database failed",
		//		zap.String("file", link.Name),
		//		zap.Error(err),
		//	)
		//} else {
		//	logger.IpfsLogger.Info("ipfs detail saved to database",
		//		zap.String("file", link.Name),
		//		zap.Float64("total_distance_km", summary.TotalDistanceKm),
		//		zap.Int("point_count", summary.PointCount),
		//	)
		//}

		list = append(list, dir)
	}
}

func (s *Service) CalcDir(ctx context.Context, rootDir string, date string) (float64, error) {
	//	解析日期
	now := carbon.Parse(date, carbon.Shanghai).StartOfDay()

	onlyDate := now.Format(carbon.DateFormat)

	split := strings.Split(onlyDate, "-")
	if len(split) < 3 {
		return 0, fmt.Errorf("日期格式错误 %v", date)
	}

	year := split[0][2:]
	month := split[1]
	day := split[2]

	fullDir := fmt.Sprintf("%s/%s/%s/%s", rootDir, year, month, day)

	deviceCodes, err := s.client.FilesLs(s.session, fullDir)
	if err != nil {
		logger.IpfsLogger.Error("device_code ipfs ls failed", zap.String("full_dir", fullDir), zap.Error(err))
		return 0, err
	}

	startTime := now.Format(carbon.DateTimeFormat)
	endTime := now.AddDay().Format(carbon.DateTimeFormat)

	var totalTurnover = 0.0

	for i, deviceCode := range deviceCodes {
		// 记录剩余未处理设备数量
		logger.IpfsLogger.Info("开始处理设备",
			zap.Int("index", i),
			zap.String("device_code", deviceCode.Name))
		//	每台设备的存储路径
		newFullPath := fmt.Sprintf("%s/%s", fullDir, deviceCode.Name)

		//	读取设备路径下单日所有文件
		gpsFiles, err := s.client.FilesLs(s.session, newFullPath)
		if err != nil {
			logger.IpfsLogger.Error("gps ipfs ls failed", zap.String("full_dir", newFullPath), zap.Error(err))
			continue
		}

		//	查询某辆车某天所有的图片地址
		var cvres []*BusImageDetailCv

		//	查询图片对应的识别乘客人数
		err = s.remoteDB.WithContext(context.Background()).
			Table("bus_image_detail_cv").
			Where("device_code = ? and collection_time >= ? and collection_time < ?", deviceCode.Name, startTime, endTime).
			Find(&cvres).Error
		if err != nil {
			logger.IpfsLogger.Error("query bus_image_detail_cv failed",
				zap.String("device_code", deviceCode.Name),
				zap.Any("start_time", startTime),
				zap.Any("end_time", endTime),
				zap.Error(err))
			continue
		}

		cvPassengers := make(map[string]int64)
		for _, cv := range cvres {
			t := cv.CollectionTime.ToTime().Format("20060102150405")
			cvPassengers[t] = cv.BaiduResult
		}

		//	解析每个文件，计算里程与周转量
		deviceTurnover := 0.0
		for _, gpsFile := range gpsFiles {
			//	仅解析 gps 文件，gps_xxxx.txt
			if !strings.HasPrefix(gpsFile.Name, "gps_") || !strings.HasSuffix(gpsFile.Name, ".txt") {
				continue
			}

			newNewFullPath := fmt.Sprintf("%s/%s", newFullPath, gpsFile.Name)

			st := time.Now()
			localPath := "./tempfile/" + gpsFile.Name
			err = s.SaveFileToLocal(newNewFullPath, localPath)
			if err != nil {
				logger.IpfsLogger.Error("save file to local failed", zap.String("file", gpsFile.Name), zap.Error(err))
				continue
			}
			logger.IpfsLogger.Info("download file done", zap.String("file", gpsFile.Name), zap.Duration("cost", time.Since(st)))

			records, err := parseFile(localPath)
			if err != nil {
				logger.IpfsLogger.Error("parse file failed", zap.String("file", gpsFile.Name), zap.Error(err))
				continue
			}

			//	删除本地临时文件
			err = os.Remove(localPath)
			if err != nil {
				logger.IpfsLogger.Error("remove local file failed", zap.String("file", gpsFile.Name), zap.Error(err))
			}

			logger.IpfsLogger.Info("parse file", zap.String("file", gpsFile.Name), zap.Int("count", len(records)))

			// 计算里程
			calculator := NewDistanceCalculator()
			summary := calculator.CalculateSummary(records)

			logger.IpfsLogger.Info("distance calculation",
				zap.String("file", gpsFile.Name),
				zap.Float64("total_distance_m", summary.TotalDistance),
				zap.Float64("total_distance_km", summary.TotalDistanceKm),
				zap.Int("point_count", summary.PointCount),
				zap.Float64("avg_speed_kmh", summary.AverageSpeed),
			)

			// 从文件名中解析出时间戳，用于计算周转量
			ext := filepath.Ext(gpsFile.Name)

			t := strings.ReplaceAll(gpsFile.Name, "gps_", "")
			t = strings.ReplaceAll(t, ext, "")

			// 查询对应的客流量，计算周转量
			if v, ok := cvPassengers[t]; ok {
				tmpTurnover := cast.ToFloat64(v) * summary.TotalDistanceKm // 周转量 = 里程 * 乘客数
				deviceTurnover += tmpTurnover
			}
		}

		totalTurnover += deviceTurnover
	}

	//	XX 日总周转量
	logger.IpfsLogger.Info(fmt.Sprintf("%s, 总周转量为：%.4f", date, totalTurnover))

	//	创建碳报告日报
	_, err = s.carbonReportDayAppService.CreateCarbonReportDay(ctx, carbonreportday_application.CreateCarbonReportDayCommand{
		Turnover:       totalTurnover,
		Baseline:       0,
		CollectionDate: timeutil.Now(now.StdTime()),
	})
	if err != nil {
		logger.IpfsLogger.Error("create carbon report day failed",
			zap.String("date", date),
			zap.Error(err),
		)
	}

	return totalTurnover, nil
}

func (s *Service) CalcDirTest(ctx context.Context, rootDir string, date string) (float64, error) {
	//	解析日期
	now := carbon.Parse(date, carbon.Shanghai).StartOfDay()

	onlyDate := now.Format(carbon.DateFormat)

	split := strings.Split(onlyDate, "-")
	if len(split) < 3 {
		return 0, fmt.Errorf("日期格式错误 %v", date)
	}

	year := split[0][2:]
	month := split[1]
	day := split[2]

	fullDir := fmt.Sprintf("%s/%s/%s/%s", rootDir, year, month, day)

	deviceCodes, err := s.client.FilesLs(s.session, fullDir)
	if err != nil {
		logger.IpfsLogger.Error("device_code ipfs ls failed", zap.String("full_dir", fullDir), zap.Error(err))
		return 0, err
	}

	startTime := now.Format(carbon.DateTimeFormat)
	endTime := now.AddDay().Format(carbon.DateTimeFormat)

	var totalTurnover = 0.0

	for i, deviceCode := range deviceCodes {
		// 记录剩余未处理设备数量
		logger.IpfsLogger.Info("开始处理设备",
			zap.Int("index", i),
			zap.String("device_code", deviceCode.Name))
		//	每台设备的存储路径
		newFullPath := fmt.Sprintf("%s/%s", fullDir, deviceCode.Name)

		//	读取设备路径下单日所有文件
		gpsFiles, err := s.client.FilesLs(s.session, newFullPath)
		if err != nil {
			logger.IpfsLogger.Error("gps ipfs ls failed", zap.String("full_dir", newFullPath), zap.Error(err))
			continue
		}

		//	查询某辆车某天所有的图片地址
		var cvres []*BusImageDetailCv

		//	查询图片对应的识别乘客人数
		err = s.remoteDB.WithContext(context.Background()).
			Table("bus_image_detail_cv").
			Where("device_code = ? and collection_time >= ? and collection_time < ?", deviceCode.Name, startTime, endTime).
			Find(&cvres).Error
		if err != nil {
			logger.IpfsLogger.Error("query bus_image_detail_cv failed",
				zap.String("device_code", deviceCode.Name),
				zap.Any("start_time", startTime),
				zap.Any("end_time", endTime),
				zap.Error(err))
			continue
		}

		cvPassengers := make(map[string]int64)
		for _, cv := range cvres {
			t := cv.CollectionTime.ToTime().Format("20060102150405")
			cvPassengers[t] = cv.BaiduResult
		}

		//	解析每个文件，计算里程与周转量
		deviceTurnover := 0.0
		for _, gpsFile := range gpsFiles {
			//	仅解析 gps 文件，gps_xxxx.txt
			if !strings.HasPrefix(gpsFile.Name, "gps_") || !strings.HasSuffix(gpsFile.Name, ".txt") {
				continue
			}

			newNewFullPath := fmt.Sprintf("%s/%s", newFullPath, gpsFile.Name)

			st := time.Now()
			localPath := "./tempfile/" + gpsFile.Name
			err = s.SaveFileToLocal(newNewFullPath, localPath)
			if err != nil {
				logger.IpfsLogger.Error("save file to local failed", zap.String("file", gpsFile.Name), zap.Error(err))
				continue
			}
			logger.IpfsLogger.Info("download file done", zap.String("file", gpsFile.Name), zap.Duration("cost", time.Since(st)))

			//	删除本地临时文件
			//err = os.Remove(localPath)
			//if err != nil {
			//	logger.IpfsLogger.Error("remove local file failed", zap.String("file", gpsFile.Name), zap.Error(err))
			//}

		}

		totalTurnover += deviceTurnover
	}

	return totalTurnover, nil
}
