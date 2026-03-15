package response

import (
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(200, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context, message string) {
	c.JSON(400, Response{
		Code:    400,
		Message: message,
	})
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message string) {
	c.JSON(401, Response{
		Code:    401,
		Message: message,
	})
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, message string) {
	c.JSON(403, Response{
		Code:    403,
		Message: message,
	})
}

// NotFound 资源不存在
func NotFound(c *gin.Context, message string) {
	c.JSON(404, Response{
		Code:    404,
		Message: message,
	})
}

// InternalError 服务器内部错误
func InternalError(c *gin.Context, message string) {
	c.JSON(500, Response{
		Code:    500,
		Message: message,
	})
}
