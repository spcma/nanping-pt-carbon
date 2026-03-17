package ipfs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Record struct {
	Timestamp time.Time
	Lat       float64
	Lon       float64
	Value     float64
}

// 解析单行，返回 Record 和错误
func parseLine(line string) (Record, error) {
	parts := strings.Split(line, "\t")
	if len(parts) != 4 {
		return Record{}, fmt.Errorf("字段数错误: 期望4, 得到%d", len(parts))
	}

	t, err := time.Parse("20060102150405", parts[0])
	if err != nil {
		return Record{}, fmt.Errorf("时间解析失败: %v", err)
	}

	lat, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return Record{}, fmt.Errorf("纬度解析失败: %v", err)
	}

	lon, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return Record{}, fmt.Errorf("经度解析失败: %v", err)
	}

	val, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return Record{}, fmt.Errorf("数值解析失败: %v", err)
	}

	return Record{Timestamp: t, Lat: lat, Lon: lon, Value: val}, nil
}

// 从文件读取并解析，跳过乱码行
func parseFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var records []Record
	scanner := bufio.NewScanner(file)

	// 可选：如果文件中有超长乱码行（超过默认64KB），可以增加缓冲区
	const maxLineSize = 1024 * 1024 // 1MB
	buf := make([]byte, maxLineSize)
	scanner.Buffer(buf, maxLineSize)

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		rec, err := parseLine(line)
		if err != nil {
			// 乱码行解析失败，跳过，可选择性输出警告（调试时可用）
			// fmt.Printf("警告: 第%d行跳过: %v\n", lineNum, err)
			continue
		}
		records = append(records, rec)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件出错: %v", err)
	}
	return records, nil
}
