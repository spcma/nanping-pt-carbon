package http

import (
	"app/internal/module/iam/application"
	"app/internal/module/iam/domain"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
	"app/internal/shared/token"
	"errors"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	appService   *application.UserService
	tokenManager token.Manager
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(appService *application.UserService, tokenManager token.Manager) *AuthHandler {
	return &AuthHandler{
		appService:   appService,
		tokenManager: tokenManager,
	}
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var cmd application.RegisterCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		logger.Warn("iam", "register - invalid request",
			zap.String("error", err.Error()),
		)
		response.BadRequest(c, "请求参数错误")
		return
	}

	id, err := h.appService.Register(platform_http.Ctx(c), cmd)
	if err != nil {
		logger.RuntimeL.Error("register failed",
			zap.String("username", cmd.Username),
			zap.Error(err),
		)

		response.InternalError(c, "注册失败")
		return
	}

	logger.Info("iam", "user registered successfully",
		zap.Int64("user_id", id),
		zap.String("username", cmd.Username),
	)
	response.Success(c, gin.H{"id": id})
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var cmd application.LoginCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		logger.Warn("iam", "login - invalid request",
			zap.String("error", err.Error()),
		)
		response.BadRequest(c, "请求参数错误")
		return
	}

	resp, err := h.appService.Login(platform_http.Ctx(c), cmd, h.tokenManager)
	if err != nil {
		logger.RuntimeL.Warn("login failed",
			zap.String("username", cmd.Username),
			zap.Error(err),
		)

		// 根据错误类型返回不同的状态码
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			response.Unauthorized(c, "用户不存在")
		case errors.Is(err, domain.ErrUserFrozen):
			response.Forbidden(c, "账号已被冻结")
		case errors.Is(err, domain.ErrInvalidPassword):
			response.Unauthorized(c, "密码错误")
		default:
			response.InternalError(c, "登录失败")
		}
		return
	}

	logger.Info("iam", "user logged in successfully",
		zap.String("username", cmd.Username),
		zap.Int64("user_id", resp.UserInfo.ID),
	)
	response.Success(c, resp)
}
