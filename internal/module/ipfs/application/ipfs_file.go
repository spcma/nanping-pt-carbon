package application

import (
	"app/internal/shared/logger"
	"app/internal/shared/timeutil"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

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

	logger.IpfsL.Info("distance result saved",
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
		logger.IpfsL.Warn("ipfs detail already exists, skip saving",
			zap.String("file", fileName),
			zap.Int64("existing_id", existingDetail.Id),
		)
		return nil
	}

	// 保存到数据库
	_, err := s.ipfsDetailAppService.CreateIpfsDetail(context.Background(), cmd)
	return err
}

// SaveContent 保存内容到文件
func (s *Service) SaveContent(ctx context.Context, content, fsDir, filename string) (string, error) {
	return s.SaveContentForClient(s.defaultClientName, ctx, content, fsDir, filename)
}

// SaveContentForClient 为指定客户端保存内容到文件
func (s *Service) SaveContentForClient(clientName string, ctx context.Context, content, fsDir, filename string) (string, error) {
	client, err := s.getClient(clientName)
	if err != nil {
		return "", err
	}

	// 打开临时文件
	fsid, err := client.client.MFOpenTempFile(client.session)
	if err != nil {
		return "", err
	}

	// 写入数据
	_, err = client.client.MFSetData(fsid, []byte(content), 0)
	if err != nil {
		return "", err
	}

	// 确保目录存在
	ok, err := s.MustDirExistsForClient(clientName, fsDir, true)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", errors.New("dir not exists & create failed")
	}

	// 保存到 NPFS
	nodePath := fsDir + "/" + filename
	ipfsid, err := client.client.MFTemp2Files(fsid, nodePath)
	if err != nil {
		return "", err
	}

	return ipfsid, nil
}

// ReadFile 读取文件
func (s *Service) ReadFile(ctx context.Context, path string) ([]byte, error) {
	return s.ReadFileForClient(ctx, s.defaultClientName, path)
}

// ReadFileForClient 为指定客户端读取文件
func (s *Service) ReadFileForClient(ctx context.Context, clientName string, path string) ([]byte, error) {

	logger.IpfsL.Info("read file", zap.String("client", clientName), zap.String("path", path))

	ext := filepath.Ext(path)

	err := s.SaveFileToLocalForClient(clientName, path, "./tmp"+ext)
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
		logger.IpfsL.Info("read file", zap.Any("record", record))
	}

	logger.IpfsL.Info("readFile", zap.Any("summary", summary))

	return nil, nil
}

func (s *Service) SaveContentToIpfs(content, fsDir, filename string) (string, error) {
	return s.SaveContentToIpfsForClient(s.defaultClientName, content, fsDir, filename)
}

// SaveContentToIpfsForClient 为指定客户端保存内容到IPFS
func (s *Service) SaveContentToIpfsForClient(clientName string, content, fsDir, filename string) (string, error) {

	logger.IpfsL.Info("save content to file", zap.String("client", clientName), zap.String("fsDir", fsDir), zap.String("filename", filename), zap.String("content", content))

	client, err := s.getClient(clientName)
	if err != nil {
		return "", err
	}

	// 打开临时文件
	fsid, err := client.client.MFOpenTempFile(client.session)
	if err != nil {
		return "", err
	}

	// 写入数据
	_, err = client.client.MFSetData(fsid, []byte(content), 0)
	if err != nil {
		return "", err
	}

	// 确保目录存在
	ok, err := s.MustDirExistsForClient(clientName, fsDir, true)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", errors.New("dir not exists & create failed")
	}

	// 将临时文件写入到IPFS最终存档
	nodePath := fsDir + "/" + filename
	ipfsid, err := client.client.MFTemp2Files(fsid, nodePath)
	if err != nil {
		return "", err
	}

	return ipfsid, nil
}

// SaveFileToIpfs 将本地文件保存到 IPFS
func (s *Service) SaveFileToIpfs(localPath, fsDir, filename string) (string, error) {
	return s.SaveFileToIpfsForClient(s.defaultClientName, localPath, fsDir, filename)
}

// SaveFileToIpfsForClient 为指定客户端将本地文件保存到 IPFS
func (s *Service) SaveFileToIpfsForClient(clientName string, localPath, fsDir, filename string) (string, error) {
	client, err := s.getClient(clientName)
	if err != nil {
		return "", err
	}

	// 打开临时文件
	fsid, err := client.client.MFOpenTempFile(client.session)
	if err != nil {
		return "", err
	}

	// 读取本地文件
	data, err := os.ReadFile(localPath)
	if err != nil {
		return "", err
	}

	// 写入数据
	_, err = client.client.MFSetData(fsid, data, 0)
	if err != nil {
		return "", err
	}

	// 确保目录存在
	ok, err := s.MustDirExistsForClient(clientName, fsDir, true)
	if err != nil {
		return "", err
	}

	if !ok {
		return "", errors.New("dir not exists & create failed")
	}

	// 将临时文件写入到IPFS最终存档
	nodePath := fsDir + "/" + filename
	ipfsid, err := client.client.MFTemp2Files(fsid, nodePath)
	if err != nil {
		return "", err
	}

	return ipfsid, nil
}

// MustDirExists 确保目录存在
func (s *Service) MustDirExists(path string, recursive bool) (bool, error) {
	return s.MustDirExistsForClient(s.defaultClientName, path, recursive)
}

// MustDirExistsForClient 为指定客户端确保目录存在
func (s *Service) MustDirExistsForClient(clientName string, path string, recursive bool) (bool, error) {
	if !s.CheckDirForClient(clientName, path) {
		err := s.CreateDirForClient(clientName, path)
		if err != nil {
			logger.IpfsL.Error("create dir failed", zap.String("client", clientName), zap.String("dir", path), zap.Error(err))
			return false, err
		}

		return true, nil
	}

	return true, nil
}

func (s *Service) Remove() {
	s.RemoveForClient(s.defaultClientName)
}

// RemoveForClient 为指定客户端删除文件
func (s *Service) RemoveForClient(clientName string) {
	client, err := s.getClient(clientName)
	if err != nil {
		logger.IpfsL.Error("get client failed", zap.String("client", clientName), zap.Error(err))
		return
	}

	//	recursive 递归 flush 直接删除
	err = client.client.FilesRm(client.session, "/tmpp/26/03/14/20260324.txt", true, true)
	if err != nil {
		return
	}
}

// ==================== 文件读取相关 ====================

// SaveFileToLocal 将 IPFS 文件保存到本地
// filePath: IPFS 文件路径
// localPath: 本地保存路径
func (s *Service) SaveFileToLocal(clientName, filePath, localPath string) error {
	return s.SaveFileToLocalForClient(clientName, filePath, localPath)
}

// SaveFileToLocalForClient 为指定客户端将 IPFS 文件保存到本地
func (s *Service) SaveFileToLocalForClient(clientName string, filePath, localPath string) error {

	logger.IpfsL.Info("save file to local", zap.String("client", clientName), zap.String("file", filePath))

	data, _, err := s.ReadFileFromIpfsForClient(clientName, filePath)
	if err != nil {
		logger.IpfsL.Error("read file from ipfs failed", zap.String("client", clientName), zap.String("file", filePath), zap.Error(err))
		return err
	}

	//_ = data

	err = os.WriteFile(localPath, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// ReadFileFromIpfs 从 Ipfs 读取文件数据
// filePath: Ipfs 文件路径（如：/np_storage/1.jpg）
// data: 文件数据，size: 文件大小，err: 错误信息
func (s *Service) ReadFileFromIpfs(filePath string) ([]byte, int64, error) {
	return s.ReadFileFromIpfsForClient(s.defaultClientName, filePath)
}

// ReadFileFromIpfsForClient 为指定客户端从 Ipfs 读取文件数据
func (s *Service) ReadFileFromIpfsForClient(clientName string, filePath string) ([]byte, int64, error) {
	client, err := s.getClient(clientName)
	if err != nil {
		return nil, 0, err
	}

	// 打开文件 URL
	fsid, err := client.client.MMOpenUrl(client.session, filePath)
	if err != nil {
		return nil, 0, err
	}
	defer client.client.MMClose(fsid)

	// 获取文件大小
	//size, err := s.client.MFGetSize(fsid)
	//if err != nil {
	//	return nil, 0, err
	//}

	// 读取文件数据
	data, err := client.client.MFGetData(fsid, 0, -1)
	if err != nil {
		return nil, 0, err
	}
	return data, int64(len(data)), nil
}

func FileNotExist(err error) bool {
	if strings.Contains(err.Error(), "no link named") {
		return true
	}
	return false
}
