package application

import (
	"app/internal/config"
	carbonreportday_application "app/internal/module/carbonreportday"
	"app/internal/module/ipfs/infrastructure"
	"app/internal/module/ipfs/rpc"
	"app/internal/shared/db"
	"app/internal/shared/entity"
	"app/internal/shared/logger"
	"app/internal/shared/timeutil"
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

const baseline = 0.11564

// 全局服务实例，供定时任务使用
var defaultService *Service

// DefaultService 获取默认的碳报告月报服务实例
func DefaultService() *Service {
	return defaultService
}

func setDefaultService(service *Service) {
	defaultService = service
}

// IpfsClientConfig IPFS客户端配置
type IpfsClientConfig struct {
	Name     string // 客户端名称标识
	Port     int    // IPFS节点端口
	BasePath string // 基础目录路径
	Enabled  bool   // 是否启用
}

// IpfsClientInstance IPFS客户端实例
type IpfsClientInstance struct {
	config           *IpfsClientConfig
	ppt              string    // 通行证
	ppt_expire_time  time.Time // 过期时间
	session          string    // 会话token
	client           *rpc.LApiStub
	guardianStop     chan struct{} // 客户端守护程序停止信号
	reconnectTrigger chan struct{} // 重连触发信号
	mutex            sync.Mutex    // 客户端操作互斥锁
}

// Service IPFS 服务
type Service struct {
	clients                   map[string]*IpfsClientInstance // 多客户端实例，key为客户端名称
	defaultClientName         string                         // 默认客户端名称
	ipfsDetailAppService      *IpfsDetailAppService
	carbonReportDayAppService *carbonreportday_application.CarbonReportDayService
}

// NewService 创建 IPFS 服务
func NewService() *Service {
	// 初始化仓储和應用服務
	ipfsDetailRepo := infrastructure.NewIpfsDetailRepository()
	ipfsDetailAppService := NewIpfsDetailAppService(ipfsDetailRepo)

	// 初始化碳报告应用服务
	carbonReportDayRepo := carbonreportday_application.NewCarbonReportDayRepository()
	carbonReportDayAppService := carbonreportday_application.NewCarbonReportDayService(carbonReportDayRepo)

	s := &Service{
		clients:                   make(map[string]*IpfsClientInstance),
		ipfsDetailAppService:      ipfsDetailAppService,
		carbonReportDayAppService: carbonReportDayAppService,
	}

	setDefaultService(s)

	if config.GlobalConfig.Ipfs.Status {
		config1 := &IpfsClientConfig{
			Name:     "4080",
			Port:     4080,
			BasePath: "/aibk",
			Enabled:  true,
		}
		err := s.InitClient(config1)
		if err != nil {
			panic(fmt.Sprintf("init config[%d] ipfs client failed: %v", 4080, err))
		}

		config2 := &IpfsClientConfig{
			Name:     "4800",
			Port:     4800,
			BasePath: "/npbus",
			Enabled:  true,
		}
		err = s.InitClient(config2)
		if err != nil {
			panic(fmt.Sprintf("init config[%d] ipfs client failed: %v", 4800, err))
		}

		//	设置默认客户端名称
		s.defaultClientName = "4080"
	}

	return s
}

var sessionInvalid = "2:session id is not safe!"

// AuthForClient 对指定客户端进行认证
func (s *Service) AuthForClient(clientName ...string) error {
	// 如果没有传入客户端名称，使用默认客户端
	name := s.defaultClientName
	if len(clientName) > 0 && clientName[0] != "" {
		name = clientName[0]
	}

	client, err := s.getClient(name)
	if err != nil {
		return err
	}

	if client.ppt == "" {
		logger.IpfsL.Error("IPFS passport is empty", zap.String("client", name))
		return errors.New("IPFS passport is empty")
	}

	// 获取会话通信证, 1h 内没有操作则默认失效
	loginReply, err := client.client.LoginWithPPT(client.ppt)
	if err != nil {
		logger.IpfsL.Error("IPFS login failed", zap.String("client", name), zap.Error(err))
		return err
	}

	client.session = loginReply.Sid

	return nil
}

// 默认配置检查目录是否存在
func (s *Service) CheckDir(dir string) bool {
	return s.CheckDirForClient(s.defaultClientName, dir)
}

// CheckDirForClient 检查指定客户端的目录
func (s *Service) CheckDirForClient(clientName string, dir string) bool {
	err := s.AuthForClient(clientName)
	if err != nil {
		return false
	}

	client, _ := s.getClient(clientName)
	stat, err := client.client.FilesStat(client.session, dir)
	if err != nil {
		if FileNotExist(err) {
			logger.IpfsL.Warn("目录不存在", zap.String("client", clientName), zap.String("dir", dir))
			return false
		}
		logger.IpfsL.Error("发生了错误", zap.String("client", clientName), zap.Error(err))
		return false
	}

	logger.IpfsL.Info("检查目录",
		zap.String("client", clientName),
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

func (s *Service) FileStat(dir string) (any, error) {
	return s.FileStatForClient(s.defaultClientName, dir)
}

// FileStatForClient 获取指定客户端的文件状态
func (s *Service) FileStatForClient(clientName string, dir string) (any, error) {
	client, err := s.getClient(clientName)
	if err != nil {
		return nil, err
	}

	stat, err := client.client.FilesStat(client.session, dir)
	if err != nil {
		if FileNotExist(err) {
			logger.IpfsL.Warn("目录不存在", zap.String("client", clientName), zap.String("dir", dir))
			return nil, err
		}
		logger.IpfsL.Error("检查目录发生了错误", zap.String("client", clientName), zap.Error(err))
		return nil, err
	}

	logger.IpfsL.Info("检查目录",
		zap.String("client", clientName),
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

	return stat, nil
}

// CreateDir 创建目录
func (s *Service) CreateDir(dir string) error {
	return s.CreateDirForClient(s.defaultClientName, dir)
}

// CreateDirForClient 为指定客户端创建目录
func (s *Service) CreateDirForClient(clientName string, dir string) error {
	if dir == "" {
		return errors.New("目录为空")
	}

	err := s.AuthForClient(clientName)
	if err != nil {
		return err
	}

	client, _ := s.getClient(clientName)
	err = client.client.FilesMkdir(client.session, dir, true)
	if err != nil {
		logger.IpfsL.Error("创建目录失败", zap.String("client", clientName), zap.String("dir", dir), zap.Error(err))
		return err
	}

	logger.IpfsL.Info("创建目录成功", zap.String("client", clientName), zap.String("dir", dir))

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
func (s *Service) ListDir(ctx context.Context, clientName, dir string) ([]*DirDto, error) {
	if clientName == "" {
		clientName = s.defaultClientName
	}
	return s.ListDirForClient(ctx, clientName, dir)
}

// ListDirForClient 列出指定客户端的目录
func (s *Service) ListDirForClient(ctx context.Context, clientName string, dir string) ([]*DirDto, error) {
	if dir == "" {
		return nil, errors.New("目录为空")
	}

	err := s.AuthForClient(clientName)
	if err != nil {
		return nil, err
	}

	client, _ := s.getClient(clientName)
	lsLinks, err := client.client.FilesLs(client.session, dir)
	if err != nil {
		logger.IpfsL.Error("列出目录失败", zap.String("client", clientName), zap.String("dir", dir), zap.Error(err))
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
	return s.ScanDirForClient(s.defaultClientName, ctx, rootDir)
}

// ScanDirForClient 为指定客户端递归扫描目录
func (s *Service) ScanDirForClient(clientName string, ctx context.Context, rootDir string) (*ScanDirResponse, error) {
	if rootDir == "" {
		return nil, errors.New("目录为空")
	}

	startTime := time.Now()
	res := &ScanDirResponse{
		RootPath: rootDir,
		Files:    make([]*ScanFileResult, 0),
	}

	logger.IpfsL.Info("正在扫描目录", zap.String("client", clientName), zap.String("rootDir", rootDir))

	err := s.AuthForClient(clientName)
	if err != nil {
		return nil, err
	}

	// 递归扫描目录
	err = s.scanDirRecursiveForClient(clientName, ctx, rootDir, res, 0, rootDir)
	if err != nil {
		logger.IpfsL.Error("扫描目录失败", zap.String("rootDir", rootDir), zap.Error(err))
		return nil, err
	}

	res.DurationMs = time.Since(startTime).Milliseconds()

	logger.IpfsL.Info("扫描目录完成",
		zap.String("client", clientName),
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
	return s.scanDirRecursiveForClient(s.defaultClientName, ctx, dir, response, depth, parentPath)
}

// scanDirRecursiveForClient 为指定客户端递归扫描目录
func (s *Service) scanDirRecursiveForClient(clientName string, ctx context.Context, dir string, response *ScanDirResponse, depth int, parentPath string) error {
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

	logger.IpfsL.Info("正在扫描目录", zap.String("client", clientName), zap.String("dir", dir), zap.Int("depth", depth))

	client, err := s.getClient(clientName)
	if err != nil {
		return err
	}

	lsLinks, err := client.client.FilesLs(client.session, dir)
	if err != nil {
		logger.IpfsL.Error("列出目录失败", zap.String("client", clientName), zap.String("dir", dir), zap.Error(err))
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
			err = s.SaveFileToLocalForClient(clientName, fullPath, "./tempfile/"+link.Name)
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
	return s.DeleteFileForClient(s.defaultClientName, ctx, path)
}

// DeleteFileForClient 为指定客户端删除文件
func (s *Service) DeleteFileForClient(clientName string, ctx context.Context, path string) error {
	if path == "" {
		return errors.New("文件路径为空")
	}

	err := s.AuthForClient(clientName)
	if err != nil {
		return err
	}

	client, _ := s.getClient(clientName)
	// recursive 递归
	// force 强制删除
	err = client.client.FilesRm(client.session, path, true, true)
	if err != nil {
		logger.IpfsL.Error("delete file error", zap.String("client", clientName), zap.Error(err))
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
func (s *Service) CalcDir(ctx context.Context, clientName string, rootDir string, date string) (float64, error) {
	return s.CalcDirForClient(ctx, clientName, rootDir, date)
}

// CalcDirForClient 为指定客户端递归扫描目录并计算周转量
func (s *Service) CalcDirForClient(ctx context.Context, clientName string, rootDir string, date string) (float64, error) {
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

	client, err := s.getClient(clientName)
	if err != nil {
		return 0, err
	}

	statInfo, err := client.client.FilesStat(client.session, rootDir)
	if err != nil {
		logger.IpfsL.Error("获取目录信息失败", zap.String("client", clientName), zap.String("dir", rootDir), zap.Error(err))
		return 0, err
	}

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
	err = s.calcDirRecursiveForClient(clientName, ctx, fullDir, date, startTime, endTime, result, 0)
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
		Hash:           statInfo.Hash,
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
		logger.IpfsL.Error("save content to ipfs failed", zap.Float64("cost", time.Now().Sub(cst).Minutes()), zap.Error(err))
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
	return s.calcDirRecursiveForClient(s.defaultClientName, ctx, dir, date, startTime, endTime, result, depth)
}

// calcDirRecursiveForClient 为指定客户端递归扫描目录并计算周转量
func (s *Service) calcDirRecursiveForClient(clientName string, ctx context.Context, dir string, date string, startTime string, endTime string, result *CalcDirResult, depth int) error {
	// 检查 context 是否被取消
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	err := s.AuthForClient(clientName)
	if err != nil {
		return err
	}

	// 检查递归深度， 避免无限递归
	if depth > maxScanDepth {
		logger.IpfsL.Warn("达到最大递归深度", zap.String("client", clientName), zap.String("dir", dir), zap.Int("depth", depth))
		return nil
	}

	logger.IpfsL.Info("正在扫描目录", zap.String("client", clientName), zap.String("dir", dir), zap.Int("depth", depth))

	client, err := s.getClient(clientName)
	if err != nil {
		return err
	}

	lsLinks, err := client.client.FilesLs(client.session, dir)
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
			err := s.calcDirRecursiveForClient(clientName, ctx, fullPath, date, startTime, endTime, result, depth+1)
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
		s.processFilesSequential(ctx, clientName, fileTasks, startTime, endTime, result)
	}

	return nil
}

// processFilesSequential 串行处理文件任务
func (s *Service) processFilesSequential(ctx context.Context, clientName string, tasks []FileTask, startTime string, endTime string, result *CalcDirResult) {
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
		turnover, err := s.processGpsFile(ctx, clientName, task.FullPath, task.FileName, task.DeviceCode, startTime, endTime)
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
func (s *Service) processGpsFile(ctx context.Context, clientName string, fullPath string, fileName string, deviceCode string, startTime string, endTime string) (float64, error) {
	st := time.Now()
	localPath := "./tempfile/" + fileName
	err := s.SaveFileToLocal(clientName, fullPath, localPath)
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

	//	获取数据源2
	remoteDB := db.GetDBWithContext(ctx, "remote")
	err = remoteDB.WithContext(ctx).
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
	return s.ReadDirTestForClient(s.defaultClientName, path)
}

// ReadDirTestForClient 为指定客户端测试读取目录
func (s *Service) ReadDirTestForClient(clientName string, path string) error {
	client, err := s.getClient(clientName)
	if err != nil {
		return err
	}

	lsLinks, err := client.client.FilesLs(client.session, path)
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

// Close 关闭连接
func (s *Service) Close() {
	// 关闭所有客户端
	for name, client := range s.clients {
		// 停止客户端守护程序
		if client.guardianStop != nil {
			close(client.guardianStop)
			logger.IpfsL.Info("IPFS 客户端守护程序停止信号已发送", zap.String("client", name))
		}

		// 关闭重连触发通道
		if client.reconnectTrigger != nil {
			close(client.reconnectTrigger)
		}

		// 关闭客户端连接
		if client.client != nil {
			client.client.Logout(client.session, "")
			logger.IpfsL.Info("IPFS 客户端连接已关闭", zap.String("client", name))
		}
	}
}

// startPassportGuardian 启动通行证自动守护程序
// 定时检查通行证是否即将过期，并自动刷新
func (s *Service) startPassportGuardian() {
	s.startPassportGuardianForClient(s.clients[s.defaultClientName])
}

// startPassportGuardianForClient 为指定客户端启动通行证自动守护程序
func (s *Service) startPassportGuardianForClient(client *IpfsClientInstance) {
	if client == nil {
		return
	}

	// 每 30 分钟检查一次通行证状态
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	logger.IpfsL.Info("通行证自动守护程序已启动", zap.String("client", client.config.Name))

	for range ticker.C {
		// 检查通行证是否将在 1 小时内过期
		if client.ppt != "" && client.ppt_expire_time.Sub(time.Now()) < 1*time.Hour {
			logger.IpfsL.Info("通行证即将过期，正在自动刷新...", zap.String("client", client.config.Name))

			// 重新获取通行证
			err := s.refreshPassportForClient(client)
			if err != nil {
				logger.IpfsL.Error("自动刷新通行证失败", zap.String("client", client.config.Name), zap.Error(err))
				// 如果刷新失败，10 分钟后重试
				continue
			}

			logger.IpfsL.Info("通行证自动刷新成功", zap.String("client", client.config.Name))
		}
	}
}

// refreshPassport 刷新通行证
func (s *Service) refreshPassport() error {
	return s.refreshPassportForClient(s.clients[s.defaultClientName])
}

// refreshPassportForClient 为指定客户端刷新通行证
func (s *Service) refreshPassportForClient(client *IpfsClientInstance) error {
	if client == nil {
		return errors.New("客户端为空")
	}

	passport, err := rpc.GetLocalPassport(client.config.Port, 24)
	if err != nil {
		logger.IpfsL.Error("IPFS refresh passport failed", zap.String("client", client.config.Name), zap.Error(err))
		return err
	}

	client.ppt = passport
	client.ppt_expire_time = time.Now().Add(24 * time.Hour)

	logger.IpfsL.Info("通行证已刷新",
		zap.String("client", client.config.Name),
		zap.Time("expire_time", client.ppt_expire_time))
	return nil
}

// startClientGuardian 启动客户端连接守护程序
// 定时检测连接状态，发现断开自动重连
func (s *Service) startClientGuardian() {
	client := s.clients[s.defaultClientName]
	if client != nil {
		s.startClientGuardianForClient(client)
	}
}

// startClientGuardianForClient 为指定客户端启动连接守护程序
func (s *Service) startClientGuardianForClient(client *IpfsClientInstance) {
	if client == nil {
		return
	}

	s.executeReconnectForClient(client)

	// 每 5 分钟检查一次连接状态
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	logger.IpfsL.Info("IPFS 客户端守护程序已启动", zap.String("client", client.config.Name))

	for {
		select {
		case <-client.guardianStop:
			logger.IpfsL.Info("IPFS 客户端守护程序已停止", zap.String("client", client.config.Name))
			return
		case <-client.reconnectTrigger:
			// 收到重连触发信号，立即重连
			logger.IpfsL.Info("收到重连触发信号，正在重连...", zap.String("client", client.config.Name))
			s.executeReconnectForClient(client)
		case <-ticker.C:
			// 定时检查连接状态
			if client.client == nil {
				logger.IpfsL.Warn("IPFS 客户端未初始化，尝试重连...", zap.String("client", client.config.Name))
				s.executeReconnectForClient(client)
				continue
			}

			// 测试连接是否正常
			_, err := client.client.FilesLs(client.session, "/")
			if err == nil {
				// 连接正常
				continue
			}

			// 连接异常，尝试重连
			if err.Error() == sessionInvalid {
				logger.IpfsL.Warn("IPFS session 已失效，正在重连...", zap.String("client", client.config.Name))
			} else {
				logger.IpfsL.Warn("IPFS 连接异常，正在重连...",
					zap.String("client", client.config.Name),
					zap.Error(err))
			}

			s.executeReconnectForClient(client)
		}
	}
}

// executeReconnect 执行重连操作（带锁保护）
func (s *Service) executeReconnect() {
	client := s.clients[s.defaultClientName]
	if client != nil {
		s.executeReconnectForClient(client)
	}
}

// executeReconnectForClient 为指定客户端执行重连操作（带锁保护）
func (s *Service) executeReconnectForClient(client *IpfsClientInstance) {
	if client == nil {
		return
	}

	client.mutex.Lock()
	defer client.mutex.Unlock()

	if err := s.reconnectClientForClient(client); err != nil {
		logger.IpfsL.Error("IPFS 重连失败", zap.String("client", client.config.Name), zap.Error(err))
	} else {
		logger.IpfsL.Info("IPFS 重连成功", zap.String("client", client.config.Name))
	}
}

// reconnectClient 重新连接客户端
func (s *Service) reconnectClient() error {
	return s.reconnectClientForClient(s.clients[s.defaultClientName])
}

// reconnectClientForClient 为指定客户端重新连接
func (s *Service) reconnectClientForClient(client *IpfsClientInstance) error {
	if client == nil {
		return errors.New("客户端为空")
	}

	// 获取服务连接
	client.client = rpc.InitLApiStubByUrl(fmt.Sprintf("127.0.0.1:%d", client.config.Port))

	// 使用通行证登录获取新 session
	err := s.AuthForClient(client.config.Name)
	if err != nil {
		logger.IpfsL.Error("IPFS 重连登录失败", zap.String("client", client.config.Name), zap.Error(err))
		return err
	}

	logger.IpfsL.Info("IPFS 客户端重连成功", zap.String("client", client.config.Name))
	return nil
}

// TriggerReconnect 触发立即重连（供其他方法调用）
// 当其他方法检测到连接断开时，可以调用此方法立即触发重连
func (s *Service) TriggerReconnect() {
	s.TriggerReconnectForClient(s.defaultClientName)
}

// TriggerReconnectForClient 触发指定客户端立即重连
func (s *Service) TriggerReconnectForClient(clientName string) {
	client, err := s.getClient(clientName)
	if err != nil {
		logger.IpfsL.Warn("客户端不存在，无法触发重连", zap.String("client", clientName))
		return
	}

	if client.reconnectTrigger != nil {
		// 非阻塞发送，如果通道已满则忽略（避免重复触发）
		select {
		case client.reconnectTrigger <- struct{}{}:
			logger.IpfsL.Info("已发送重连触发信号", zap.String("client", clientName))
		default:
			logger.IpfsL.Debug("重连信号已在队列中，忽略重复触发", zap.String("client", clientName))
		}
	} else {
		logger.IpfsL.Warn("重连通道未初始化，无法触发重连", zap.String("client", clientName))
	}
}

// parseDirByPort 根据端口号返回不同的日期格式化路径
func parseDirByPort(port int, date time.Time) (string, error) {

	if date.IsZero() {
		date = time.Now()
	}

	// 获取年月日
	year := date.Year() % 100 // 取两位年份
	month := int(date.Month())
	day := date.Day()

	// 格式化为两位数
	yearStr := fmt.Sprintf("%02d", year)
	monthStr := fmt.Sprintf("%02d", month)
	dayStr := fmt.Sprintf("%02d", day)

	switch port {
	case 4080:
		// 格式: /aibk/26/01/01/
		return fmt.Sprintf("/aibk/%s/%s/%s/", yearStr, monthStr, dayStr), nil
	case 4800:
		// 格式: /npbk/260101/
		return fmt.Sprintf("/npbus/%s%s%s/", yearStr, monthStr, dayStr), nil
	default:
		return "", fmt.Errorf("不支持的端口号: %d", port)
	}
}

// ==================== 客户端管理方法 ====================

// InitClient 初始化一个新的IPFS客户端
func (s *Service) InitClient(config *IpfsClientConfig) error {
	if config == nil {
		return errors.New("客户端配置不能为空")
	}

	if !config.Enabled {
		logger.IpfsL.Info("客户端未启用", zap.String("name", config.Name))
		return nil
	}

	// 创建客户端实例
	client := &IpfsClientInstance{
		config: config,
	}

	// 获取通行证
	err := s.refreshPassportForClient(client)
	if err != nil {
		return fmt.Errorf("获取通行证失败: %v", err)
	}

	// 初始化RPC客户端
	client.client = rpc.InitLApiStubByUrl(fmt.Sprintf("127.0.0.1:%d", config.Port))

	// 先将客户端存储到映射中，以便后续认证时可以找到
	s.clients[config.Name] = client

	// 登录获取session（此时client已经在s.clients中）
	err = s.AuthForClient(config.Name)
	if err != nil {
		// 登录失败，从映射中移除
		delete(s.clients, config.Name)
		return fmt.Errorf("登录失败: %v", err)
	}

	// 初始化通道
	client.guardianStop = make(chan struct{})
	client.reconnectTrigger = make(chan struct{}, 1)

	// 启动守护程序
	go s.startClientGuardianForClient(client)

	logger.IpfsL.Info("IPFS客户端初始化成功",
		zap.String("name", config.Name),
		zap.Int("port", config.Port),
		zap.String("basePath", config.BasePath))

	return nil
}

// getClient 获取指定名称的客户端
func (s *Service) getClient(name string) (*IpfsClientInstance, error) {
	client, ok := s.clients[name]
	if !ok {
		return nil, fmt.Errorf("客户端不存在: %s", name)
	}
	return client, nil
}

// GetClient 公开方法：获取指定名称的客户端
func (s *Service) GetClient(name string) (*IpfsClientInstance, error) {
	return s.getClient(name)
}

// GetAllClients 获取所有客户端
func (s *Service) GetAllClients() map[string]*IpfsClientInstance {
	return s.clients
}

// SetDefaultClient 设置默认客户端
func (s *Service) SetDefaultClient(name string) error {
	if _, ok := s.clients[name]; !ok {
		return fmt.Errorf("客户端不存在: %s", name)
	}
	s.defaultClientName = name
	return nil
}

// GetDefaultClientName 获取默认客户端名称
func (s *Service) GetDefaultClientName() string {
	return s.defaultClientName
}
