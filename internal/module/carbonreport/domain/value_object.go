package domain

import (
	"app/internal/shared/valueobject"
)

// 值对象类型别名，方便使用

type (
	// FilePath 文件路径值对象（使用 shared 的封装）
	FilePath = string

	// FileName 文件名值对象
	FileName = string

	// FileContent 文件内容值对象
	FileContent = []byte

	// DirPath 目录路径值对象
	DirPath = string

	// IPFSID IPFS 文件 ID 值对象
	IPFSID = string

	// Coordinate 坐标值对象（复用 shared）
	Coordinate = valueobject.Coordinate

	// License 许可证值对象（复用 shared）
	License = valueobject.License

	// TimeRange 时间范围值对象（复用 shared）
	TimeRange = valueobject.TimeRange
)
