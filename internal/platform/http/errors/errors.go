package errors

import "fmt"

// HTTPError HTTP 错误
type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.Code, e.Message)
}

// 常见错误定义
var (
	ErrBadRequest         = &HTTPError{Code: 400, Message: "请求参数错误"}
	ErrUnauthorized       = &HTTPError{Code: 401, Message: "未授权访问"}
	ErrForbidden          = &HTTPError{Code: 403, Message: "禁止访问"}
	ErrNotFound           = &HTTPError{Code: 404, Message: "资源不存在"}
	ErrInternalServer     = &HTTPError{Code: 500, Message: "服务器内部错误"}
	ErrServiceUnavailable = &HTTPError{Code: 503, Message: "服务不可用"}
)

// NewHTTPError 创建新的 HTTP 错误
func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{
		Code:    code,
		Message: message,
	}
}

// WrapError 包装错误消息
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return &HTTPError{
		Code:    500,
		Message: fmt.Sprintf("%s: %v", message, err),
	}
}
