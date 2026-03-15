package error

// Error 领域错误
type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

// NewError 创建通用错误
func NewError(code, message string) *Error {
	return &Error{Code: code, Message: message}
}
