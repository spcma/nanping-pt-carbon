package domain

import (
	sharederror "app/internal/shared/error"
)

// 领域错误定义（使用 shared 的 Error 结构）
var (
	ErrFileNotFound      = sharederror.NewError("FILE_NOT_FOUND", "文件不存在")
	ErrDirNotFound       = sharederror.NewError("DIR_NOT_FOUND", "目录不存在")
	ErrFileAlreadyExists = sharederror.NewError("FILE_ALREADY_EXISTS", "文件已存在")
	ErrDirAlreadyExists  = sharederror.NewError("DIR_ALREADY_EXISTS", "目录已存在")
	ErrInvalidPath       = sharederror.NewError("INVALID_PATH", "无效的路径")
	ErrReadFileFailed    = sharederror.NewError("READ_FILE_FAILED", "读取文件失败")
	ErrWriteFileFailed   = sharederror.NewError("WRITE_FILE_FAILED", "写入文件失败")
	ErrDeleteFileFailed  = sharederror.NewError("DELETE_FILE_FAILED", "删除文件失败")
)
