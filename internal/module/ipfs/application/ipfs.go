package application

import (
	"app/internal/config"
	carbonreportday_application "app/internal/module/carbonreportday"
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
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/dromara/carbon/v2"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

const baseline = 0.11564

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
	carbonReportDayRepo := carbonreportday_application.NewCarbonReportDayRepository(db)
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

// CreateFsClient 创建文件系统客户端
func CreateFsClient() (*rpc.LApiStub, string, error) {
	strPPT, err := rpc.GetLocalPassport(config.GlobalConfig.Ipfs.Port, 24)
	if err != nil {
		return nil, "", err
	}

	client := rpc.InitLApiStubByUrl(fmt.Sprintf("127.0.0.1:%d", config.GlobalConfig.Ipfs.Port))

	loginReply, err := client.LoginWithPPT(strPPT)
	if err != nil {
		return nil, "", err
	}

	return client, loginReply.Sid, nil
}

func (s *Service) CheckDir(dir string) bool {
	stat, err := s.client.FilesStat(s.session, dir)
	if err != nil {
		if FileNotExist(err) {
			logger.IpfsL.Warn("目录不存在", zap.String("dir", dir))
			return false
		}
		logger.IpfsL.Error("发生了错误", zap.Error(err))
		return false
	}

	logger.IpfsL.Info("检查目录",
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
func (s *Service) CreateDir(dir string) error {
	if dir == "" {
		return errors.New("目录为空")
	}

	err := s.client.FilesMkdir(s.session, dir, true)
	if err != nil {
		logger.IpfsL.Error("创建目录失败", zap.String("dir", dir), zap.Error(err))
		return err
	}

	logger.IpfsL.Info("创建目录成功", zap.String("dir", dir))

	return nil
}

// DirDto 目录项
type DirDto struct {
	Name      string
	Hash      string
	Size      uint64
	Type      int    // 1 dir 0 file
	FileType  string // .jpg .txt
	Timestamp string // 时间戳
	Seq       int    // 序号
}

// ListDir 列出目录
//
//	gps_20260301180732.txt
func (s *Service) ListDir(ctx context.Context, dir string) ([]*DirDto, error) {
	if dir == "" {
		return nil, errors.New("目录为空")
	}

	lsLinks, err := s.client.FilesLs(s.session, dir)
	if err != nil {
		logger.IpfsL.Error("列出目录失败", zap.String("dir", dir), zap.Error(err))
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

		if strings.HasPrefix(link.Name, "gps_") && strings.HasSuffix(link.Name, ".txt") {
			//	截取文件名中的时间戳 20260301180452
			timestamp := strings.TrimSuffix(strings.TrimPrefix(link.Name, "gps_"), ".txt")
			dirDto.Timestamp = timestamp
			//	文件类型，截取文件后缀
			dirDto.FileType = filepath.Ext(link.Name)
		}

		list = append(list, dirDto)
	}

	return list, nil
}

// ScanFileResult 扫描文件结果
type ScanFileResult struct {
	Path       string `json:"path"`       // 完整路径
	Name       string `json:"name"`       // 文件名
	Size       uint64 `json:"size"`       // 文件大小
	FileType   string `json:"fileType"`   // 文件类型
	Timestamp  string `json:"timestamp"`  // 时间戳（从文件名提取）
	Depth      int    `json:"depth"`      // 目录深度
	ParentPath string `json:"parentPath"` // 父目录路径
}

// CalcDirResponse 计算目录响应
type CalcDirResponse struct {
	RootPath      string          `json:"rootPath"`      // 根目录路径
	TotalFiles    int             `json:"totalFiles"`    // 总文件数
	TotalDirs     int             `json:"totalDirs"`     // 总目录数
	TotalDistance float64         `json:"totalDistance"` // 总里程(km)
	TotalTurnover float64         `json:"totalTurnover"` // 总周转量
	DeviceResults []*DeviceResult `json:"deviceResults"` // 设备计算结果
	DurationMs    int64           `json:"durationMs"`    // 耗时（毫秒）
}

// DeviceResult 设备计算结果
type DeviceResult struct {
	DeviceCode      string  `json:"deviceCode"`      // 设备编号
	TotalDistanceKm float64 `json:"totalDistanceKm"` // 总里程(km)
	Turnover        float64 `json:"turnover"`        // 周转量
	FileCount       int     `json:"fileCount"`       // 文件数
}

// CalcFileResult 计算文件结果
type CalcFileResult struct {
	Path           string  `json:"path"`           // 完整路径
	Name           string  `json:"name"`           // 文件名
	DeviceCode     string  `json:"deviceCode"`     // 设备编号
	DistanceKm     float64 `json:"distanceKm"`     // 里程(km)
	Timestamp      string  `json:"timestamp"`      // 时间戳
	PassengerCount int64   `json:"passengerCount"` // 乘客数
	Turnover       float64 `json:"turnover"`       // 周转量
}

// ScanDirResponse 扫描目录响应
type ScanDirResponse struct {
	RootPath   string            `json:"rootPath"`   // 根目录路径
	TotalFiles int               `json:"totalFiles"` // 总文件数
	TotalDirs  int               `json:"totalDirs"`  // 总目录数
	TotalSize  uint64            `json:"totalSize"`  // 总大小
	Files      []*ScanFileResult `json:"files"`      // 文件列表
	DurationMs int64             `json:"durationMs"` // 耗时（毫秒）
}

// maxScanDepth 最大递归扫描深度
const maxScanDepth = 50

// ScanDir 递归扫描目录，遍历所有子目录和文件
func (s *Service) ScanDir(ctx context.Context, rootDir string) (*ScanDirResponse, error) {
	if rootDir == "" {
		return nil, errors.New("目录为空")
	}

	startTime := time.Now()
	res := &ScanDirResponse{
		RootPath: rootDir,
		Files:    make([]*ScanFileResult, 0),
	}

	logger.IpfsL.Info("正在扫描目录", zap.String("rootDir", rootDir))

	// 递归扫描目录
	err := s.scanDirRecursive(ctx, rootDir, res, 0, rootDir)
	if err != nil {
		logger.IpfsL.Error("扫描目录失败", zap.String("rootDir", rootDir), zap.Error(err))
		return nil, err
	}

	res.DurationMs = time.Since(startTime).Milliseconds()

	logger.IpfsL.Info("扫描目录完成",
		zap.String("rootDir", rootDir),
		zap.Int("totalFiles", res.TotalFiles),
		zap.Int("totalDirs", res.TotalDirs),
		zap.Uint64("totalSize", res.TotalSize),
		zap.Int64("durationMs", res.DurationMs),
	)

	return res, nil
}

// scanDirRecursive 递归扫描目录
func (s *Service) scanDirRecursive(ctx context.Context, dir string, response *ScanDirResponse, depth int, parentPath string) error {
	// 检查 context 是否被取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// 检查递归深度
	if depth > maxScanDepth {
		logger.IpfsL.Warn("达到最大递归深度", zap.String("dir", dir), zap.Int("depth", depth))
		return nil
	}

	logger.IpfsL.Info("正在扫描目录", zap.String("dir", dir), zap.Int("depth", depth))

	lsLinks, err := s.client.FilesLs(s.session, dir)
	if err != nil {
		logger.IpfsL.Error("列出目录失败", zap.String("dir", dir), zap.Error(err))
		return err
	}

	for _, link := range lsLinks {
		fullPath := path.Join(dir, link.Name)

		if link.Type == 1 { // 目录
			response.TotalDirs++

			// 递归扫描子目录
			err := s.scanDirRecursive(ctx, fullPath, response, depth+1, dir)
			if err != nil {
				logger.IpfsL.Warn("扫描子目录失败", zap.String("dir", fullPath), zap.Error(err))
			}
		} else if link.Type == 0 { // 文件

			//	仅解析 gps 文件，gps_xxxx.txt
			if !strings.HasPrefix(link.Name, "gps_") || !strings.HasSuffix(link.Name, ".txt") {
				continue
			}

			response.TotalFiles++
			response.TotalSize += link.Size

			st := time.Now()
			err = s.SaveFileToLocal(fullPath, "./tempfile/"+link.Name)
			if err != nil {
				logger.IpfsL.Error("保存文件失败", zap.String("file", fullPath), zap.Error(err))
				return err
			}
			logger.IpfsL.Info("download file done", zap.Duration("cost", time.Since(st)), zap.String("fullPath", fullPath), zap.Int64("size", int64(link.Size)), zap.String("file", link.Name))

			// 创建文件结果对象
			result := &ScanFileResult{
				Path:       fullPath,
				Name:       link.Name,
				Size:       link.Size,
				Depth:      depth + 1,
				ParentPath: dir,
			}

			// 提取文件信息
			if strings.HasSuffix(link.Name, ".txt") {
				result.FileType = ".txt"
				if strings.HasPrefix(link.Name, "gps_") {
					result.Timestamp = strings.TrimSuffix(strings.TrimPrefix(link.Name, "gps_"), ".txt")
				}
			} else {
				result.FileType = filepath.Ext(link.Name)
			}

			logger.IpfsL.Info("扫描文件", zap.Any("result", result))
		}
	}

	return nil
}

// DeleteFile 删除文件
func (s *Service) DeleteFile(ctx context.Context, path string) error {
	if path == "" {
		return errors.New("文件路径为空")
	}

	// recursive 递归
	// force 强制删除
	err := s.client.FilesRm(s.session, path, true, true)
	if err != nil {
		logger.IpfsL.Error("delete file error", zap.Error(err))
		return err
	}

	return nil
}

// UploadFile 上传文件
func (s *Service) UploadFile(tmpfilePath, uploadDir, filename string) (string, error) {

	//	上传文件到 IPFS
	ipfsid, err := s.SaveFileToIpfs(tmpfilePath, uploadDir, filename)
	if err != nil {
		return "", err
	}

	logger.IpfsL.Info("上传文件成功", zap.String("filename", filename), zap.String("ipfsid", ipfsid))

	return ipfsid, nil
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

	logger.IpfsL.Info("handle with dir", zap.String("dir", req.Dir),
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
			logger.IpfsL.Info("skip dir", zap.String("dir", link.Name))
			continue
		}

		if strings.Contains(link.Name, ".jpg") {
			logger.IpfsL.Info("skip image", zap.String("file", link.Name))
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

		logger.IpfsL.Info("parse file", zap.String("file", link.Name),
			zap.Int("count", len(records)),
		)

		// 计算里程并保存到数据库
		calculator := NewDistanceCalculator()
		summary := calculator.CalculateSummary(records)

		logger.IpfsL.Info("distance calculation",
			zap.String("file", link.Name),
			zap.Float64("total_distance_m", summary.TotalDistance),
			zap.Float64("total_distance_km", summary.TotalDistanceKm),
			zap.Int("point_count", summary.PointCount),
			zap.Float64("avg_speed_kmh", summary.AverageSpeed),
		)

		dirDto := DirDto{
			Name: link.Name,
			Hash: link.Hash,
			Size: link.Size,
			Type: link.Type,
		}

		if strings.HasPrefix(link.Name, "gps_") && strings.HasSuffix(link.Name, ".txt") {
			//	截取文件名中的时间戳 20260301180452
			timestamp := strings.TrimSuffix(strings.TrimPrefix(link.Name, "gps_"), ".txt")
			dirDto.Timestamp = timestamp
			//	文件类型, 截取文件后缀
			dirDto.FileType = filepath.Ext(link.Name)
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

		list = append(list, dirDto)
	}
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

// CalcDir 递归扫描目录并计算周转量
// rootDir: 要扫描的根目录（直接从此目录开始递归）
// date: 日期，格式为 "2026-03-27"，用于查询数据库和生成报告
func (s *Service) CalcDir(ctx context.Context, rootDir string, date string) (float64, error) {
	cst := time.Now()

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

	// 直接使用传入的 rootDir 作为扫描起点
	fullDir := rootDir

	startTime := now.Copy().Format(carbon.DateTimeFormat)
	endTime := now.Copy().AddDay().Format(carbon.DateTimeFormat)

	// 初始化计算结果
	result := &CalcDirResult{
		TotalTurnover: 0,
		DeviceResults: make(map[string]*DeviceCalcResult),
	}

	logger.IpfsL.Info("开始递归扫描计算",
		zap.String("rootDir", rootDir),
		zap.String("date", date),
		zap.String("fullDir", fullDir),
		zap.String("startTime", startTime),
		zap.String("endTime", endTime),
	)

	// 递归扫描目录并计算
	err := s.calcDirRecursive(ctx, fullDir, date, startTime, endTime, result, 0)
	if err != nil {
		logger.IpfsL.Error("递归计算目录失败", zap.String("full_dir", fullDir), zap.Error(err))
		return 0, err
	}

	totalTurnover := result.TotalTurnover

	//	XX 日总周转量
	logger.IpfsL.Info(fmt.Sprintf("%s, 总周转量为：%.4f", date, totalTurnover))

	//	创建碳报告日报
	_, err = s.carbonReportDayAppService.CreateCarbonReportDay(ctx, carbonreportday_application.CreateCarbonReportDayCommand{
		Turnover:       totalTurnover,
		Baseline:       0,
		CollectionDate: timeutil.Now(now.StdTime()),
	})
	if err != nil {
		logger.IpfsL.Error("create carbon report day failed",
			zap.String("date", date),
			zap.Error(err),
		)
	}

	//	 保存周转量结果到文件中
	saveDir := fmt.Sprintf("%s/%s/%s/%s", "/tmpp", year, month, day)

	filename := fmt.Sprintf("%s.txt", now.Format(carbon.ShortDateTimeFormat))

	//	计算节碳量
	value := totalTurnover * baseline

	saveContent := strings.Builder{}
	saveContent.WriteString(now.ToDateString())
	saveContent.WriteString("\t")
	saveContent.WriteString(cast.ToString(totalTurnover))
	saveContent.WriteString("\t")
	saveContent.WriteString(cast.ToString(baseline))
	saveContent.WriteString("\t")
	saveContent.WriteString(cast.ToString(value))

	ipfsid, err := s.SaveContentToIpfs(saveContent.String(), saveDir, filename)
	if err != nil {
		logger.IpfsL.Error("save content to ipfs failed", zap.Error(err))
		return 0, err
	}

	logger.IpfsL.Info("save content to ipfs done", zap.String("ipfs_id", ipfsid), zap.Float64("cost", time.Now().Sub(cst).Minutes()))

	return totalTurnover, nil
}

// CalcDirResult 计算结果汇总
type CalcDirResult struct {
	TotalTurnover float64
	DeviceResults map[string]*DeviceCalcResult // deviceCode -> result
}

// DeviceCalcResult 设备计算结果
type DeviceCalcResult struct {
	DeviceCode    string
	Turnover      float64
	FileCount     int
	TotalDistance float64
}

// FileTask 文件处理任务
type FileTask struct {
	FullPath   string
	FileName   string
	DeviceCode string
}

// FileResult 文件处理结果
type FileResult struct {
	Task     FileTask
	Turnover float64
	Err      error
}

// calcDirRecursive 递归扫描目录并计算周转量（并发模式）
func (s *Service) calcDirRecursive(ctx context.Context, dir string, date string, startTime string, endTime string, result *CalcDirResult, depth int) error {
	// 检查 context 是否被取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// 检查递归深度， 避免无限递归
	if depth > maxScanDepth {
		logger.IpfsL.Warn("达到最大递归深度", zap.String("dir", dir), zap.Int("depth", depth))
		return nil
	}

	logger.IpfsL.Info("正在扫描目录", zap.String("dir", dir), zap.Int("depth", depth))

	lsLinks, err := s.client.FilesLs(s.session, dir)
	if err != nil {
		logger.IpfsL.Error("列出目录失败", zap.String("dir", dir), zap.Error(err))
		return err
	}

	// 收集所有文件任务
	var fileTasks []FileTask

	for _, link := range lsLinks {
		fullPath := path.Join(dir, link.Name)

		switch link.Type {
		case 0: // file
			//	仅解析 gps 文件，gps_xxxx.txt
			if !strings.HasPrefix(link.Name, "gps_") || !strings.HasSuffix(link.Name, ".txt") {
				continue
			}

			// 从路径中提取设备编码（假设路径结构为 .../year/month/day/deviceCode/gps_xxx.txt）
			deviceCode := s.extractDeviceCodeFromPath(fullPath)
			if deviceCode == "" {
				logger.IpfsL.Warn("无法从路径提取设备编码", zap.String("path", fullPath))
				continue
			}

			// 初始化设备结果
			if _, ok := result.DeviceResults[deviceCode]; !ok {
				result.DeviceResults[deviceCode] = &DeviceCalcResult{
					DeviceCode: deviceCode,
				}
			}

			// 收集文件任务
			fileTasks = append(fileTasks, FileTask{
				FullPath:   fullPath,
				FileName:   link.Name,
				DeviceCode: deviceCode,
			})
		case 1: // dir
			// 尝试从目录名解析设备编码（如果目录结构符合预期）
			deviceCode := link.Name

			// 递归扫描子目录
			err := s.calcDirRecursive(ctx, fullPath, date, startTime, endTime, result, depth+1)
			if err != nil {
				logger.IpfsL.Warn("扫描子目录失败", zap.String("dir", fullPath), zap.Error(err))
			}

			// 如果子目录有计算结果，记录到设备结果中
			if deviceResult, ok := result.DeviceResults[deviceCode]; ok && deviceResult != nil {
				logger.IpfsL.Info("设备计算完成",
					zap.String("device_code", deviceCode),
					zap.Float64("turnover", deviceResult.Turnover),
					zap.Int("file_count", deviceResult.FileCount),
				)
			}
		default:
			logger.IpfsL.Warn("未知的文件类型", zap.String("dir", dir), zap.String("name", link.Name))
		}
	}

	// 串行处理文件任务
	if len(fileTasks) > 0 {
		s.processFilesSequential(ctx, fileTasks, startTime, endTime, result)
	}

	return nil
}

// processFilesSequential 串行处理文件任务
func (s *Service) processFilesSequential(ctx context.Context, tasks []FileTask, startTime string, endTime string, result *CalcDirResult) {
	taskCount := len(tasks)
	logger.IpfsL.Info("开始串行处理文件",
		zap.Int("totalTasks", taskCount),
	)

	// 串行处理每个文件任务
	var processedCount int
	for _, task := range tasks {
		// 检查 context 是否被取消
		select {
		case <-ctx.Done():
			logger.IpfsL.Error("处理 gps 文件被取消",
				zap.String("file", task.FileName),
				zap.String("device_code", task.DeviceCode),
				zap.Error(ctx.Err()),
			)
			continue
		default:
		}

		// 处理文件
		turnover, err := s.processGpsFile(ctx, task.FullPath, task.FileName, task.DeviceCode, startTime, endTime)
		processedCount++

		if err != nil {
			logger.IpfsL.Error("处理 gps 文件失败",
				zap.String("file", task.FileName),
				zap.String("device_code", task.DeviceCode),
				zap.Error(err),
			)
			continue
		}

		// 累加设备结果
		deviceCode := task.DeviceCode
		result.DeviceResults[deviceCode].Turnover += turnover
		result.DeviceResults[deviceCode].FileCount++
		result.TotalTurnover += turnover

		logger.IpfsL.Info("文件处理完成",
			zap.String("file", task.FileName),
			zap.String("device_code", deviceCode),
			zap.Float64("turnover", turnover),
			zap.Int("progress", processedCount),
			zap.Int("total", taskCount),
		)
	}

	logger.IpfsL.Info("串行处理文件完成", zap.Int("processedCount", processedCount), zap.Int("totalTasks", taskCount))
}

// extractDeviceCodeFromPath 从路径中提取设备编码
// 路径格式: /root/year/month/day/deviceCode/gps_xxx.txt
func (s *Service) extractDeviceCodeFromPath(filePath string) string {
	parts := strings.Split(filePath, "/")
	if len(parts) >= 2 {
		// 设备编码是文件所在目录的名称
		return parts[len(parts)-2]
	}
	return ""
}

// processGpsFile 处理单个 GPS 文件，返回该文件的周转量
func (s *Service) processGpsFile(ctx context.Context, fullPath string, fileName string, deviceCode string, startTime string, endTime string) (float64, error) {
	st := time.Now()
	localPath := "./tempfile/" + fileName
	err := s.SaveFileToLocal(fullPath, localPath)
	if err != nil {
		logger.IpfsL.Error("save file to local failed", zap.String("file", fileName), zap.Error(err))
		return 0, err
	}
	logger.IpfsL.Info("download file done", zap.String("file", fileName), zap.Duration("cost", time.Since(st)))

	records, err := parseFile(localPath)
	if err != nil {
		logger.IpfsL.Error("parse file failed", zap.String("file", fileName), zap.Error(err))
		return 0, err
	}

	//	删除本地临时文件
	err = os.Remove(localPath)
	if err != nil {
		logger.IpfsL.Error("remove local file failed", zap.String("file", fileName), zap.Error(err))
	}

	logger.IpfsL.Info("parse file", zap.String("file", fileName), zap.Int("count", len(records)))

	// 计算里程
	calculator := NewDistanceCalculator()
	summary := calculator.CalculateSummary(records)

	logger.IpfsL.Info("distance calculation",
		zap.String("file", fileName),
		zap.Float64("total_distance_m", summary.TotalDistance),
		zap.Float64("total_distance_km", summary.TotalDistanceKm),
		zap.Int("point_count", summary.PointCount),
		zap.Float64("avg_speed_kmh", summary.AverageSpeed),
	)

	//	解析文件名 获取时间戳
	timestamp := strings.TrimSuffix(strings.TrimPrefix(fileName, "gps_"), ".txt")

	// 查询对应的客流量
	var cvres []*BusImageDetailCv
	err = s.remoteDB.WithContext(ctx).
		Table("bus_image_detail_cv").
		Where("device_code = ? and collection_time >= ? and collection_time < ?",
			deviceCode, startTime, endTime).
		Find(&cvres).Error
	if err != nil {
		logger.IpfsL.Error("query bus_image_detail_cv failed",
			zap.String("device_code", deviceCode),
			zap.Error(err))
		// 查询失败不中断，只是没有乘客数据
		return 0, nil
	}

	// 计算周转量
	var deviceTurnover float64
	for _, cv := range cvres {
		t := cv.CollectionTime.ToTime().Format("20060102150405")
		if t == timestamp {
			tmpTurnover := cast.ToFloat64(cv.BaiduResult) * summary.TotalDistanceKm // 周转量 = 里程 * 乘客数
			deviceTurnover += tmpTurnover
		}
	}

	return deviceTurnover, nil
}

func (s *Service) ReadDirTest(path string) error {
	lsLinks, err := s.client.FilesLs(s.session, path)
	if err != nil {
		logger.IpfsL.Error("ipfs ls failed", zap.Error(err))
		return err
	}

	logger.IpfsL.Info("ipfs ls done", zap.Int("count", len(lsLinks)))

	for _, link := range lsLinks {
		if link.Type == 0 { // 1 dir 0 file

		}
	}

	return nil
}

func (s *Service) CalcDirTest2() (any, error) {

	path := "/aibk/26/03/15/xNUr48spW1gR2bQTSRURMCl_cII"

	logger.IpfsL.Info("ipfs ls", zap.String("path", path))

	files, err := s.client.FilesLs(s.session, path)
	if err != nil {
		logger.IpfsL.Error("ipfs ls failed", zap.Error(err))
		return nil, err
	}

	logger.IpfsL.Info("ipfs ls done", zap.Int("count", len(files)))

	for _, file := range files {
		st := time.Now()
		logger.IpfsL.Info("download file", zap.String("file", file.Name))
		err = s.SaveFileToLocal("/aibk/26/03/15/xNUr48spW1gR2bQTSRURMCl_cII/"+file.Name, "./tmpfile/"+file.Name)
		if err != nil {
			logger.IpfsL.Error("save file to local failed", zap.String("file", file.Name), zap.Error(err))
			continue
		}
		logger.IpfsL.Info("download file done", zap.String("file", file.Name), zap.Duration("cost", time.Since(st)))
	}

	logger.IpfsL.Info("download file done")
	return nil, nil
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

	fmt.Println("full_dir", fullDir)
	deviceCodes, err := s.client.FilesLs(s.session, fullDir)
	if err != nil {
		logger.IpfsL.Error("device_code ipfs ls failed", zap.String("full_dir", fullDir), zap.Error(err))
		return 0, err
	}

	logger.IpfsL.Info("device_code ipfs ls done", zap.Int("count", len(deviceCodes)))

	startTime := now.Format(carbon.DateTimeFormat)
	endTime := now.AddDay().Format(carbon.DateTimeFormat)

	var totalTurnover = 0.0

	for i, deviceCode := range deviceCodes {
		// 记录剩余未处理设备数量
		logger.IpfsL.Info("开始处理设备",
			zap.Int("index", i),
			zap.String("device_code", deviceCode.Name))
		//	每台设备的存储路径
		newFullPath := fmt.Sprintf("%s/%s", fullDir, deviceCode.Name)

		//	读取设备路径下单日所有文件
		gpsFiles, err := s.client.FilesLs(s.session, newFullPath)
		if err != nil {
			logger.IpfsL.Error("gps ipfs ls failed", zap.String("full_dir", newFullPath), zap.Error(err))
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
			logger.IpfsL.Error("query bus_image_detail_cv failed",
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
				logger.IpfsL.Error("save file to local failed", zap.String("file", gpsFile.Name), zap.Error(err))
				continue
			}
			logger.IpfsL.Info("download file done", zap.String("file", gpsFile.Name), zap.Duration("cost", time.Since(st)))

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

// Close 关闭连接
func (s *Service) Close() {
	if s.client != nil {
		s.client.Logout(s.session, "")
	}
}
