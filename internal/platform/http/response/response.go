package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应 200
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应 200
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应 200
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequest 请求参数错误 400
func BadRequest(c *gin.Context, message string) {
	if message == "" {
		message = "bad request"
	}

	c.JSON(400, Response{
		Code:    400,
		Message: message,
	})
}

// Unauthorized 未授权 401
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "unauthorized"
	}

	c.JSON(401, Response{
		Code:    401,
		Message: message,
	})
}

// Forbidden 禁止访问 403
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "forbidden"
	}

	c.JSON(403, Response{
		Code:    403,
		Message: message,
	})
}

// NotFound 资源不存在 404
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "resource not found"
	}

	c.JSON(404, Response{
		Code:    404,
		Message: message,
	})
}

// InternalError 服务器内部错误 500
func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = "internal server error"
	}

	c.JSON(500, Response{
		Code:    500,
		Message: message,
	})
}

// BadGateway 网关错误 502
func BadGateway(c *gin.Context, message string) {
	if message == "" {
		message = "bad gateway"
	}

	c.JSON(502, Response{
		Code:    502,
		Message: message,
	})
}

// ServiceUnavailable 服务不可用 503
func ServiceUnavailable(c *gin.Context, message string) {
	if message == "" {
		message = "service unavailable"
	}

	c.JSON(503, Response{
		Code:    503,
		Message: message,
	})
}

// GatewayTimeout 网关超时 504
func Timeout(c *gin.Context, message string) {
	if message == "" {
		message = "gateway timeout"
	}

	c.JSON(504, Response{
		Code:    504,
		Message: message,
	})
}
