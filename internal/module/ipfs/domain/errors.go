package domain

import "errors"

// IPFS 领域错误定义
var (
	ErrIpfsDetailNotFound      = errors.New("IPFS 详情不存在")
	ErrIpfsDetailAlreadyExists = errors.New("IPFS 详情已存在")
	ErrInvalidDeviceCode       = errors.New("无效的设备编码")
	ErrInvalidFilename         = errors.New("无效的文件名")
)
