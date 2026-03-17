package http

import (
	"app/internal/module/iam/application"
	"app/internal/module/iam/domain"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
	"app/internal/shared/token"
	"net/http"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	appService   *application.SysUserAppService
	tokenManager token.Manager
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(appService *application.SysUserAppService, tokenManager token.Manager) *AuthHandler {
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
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.appService.Register(platform_http.Ctx(c), cmd)
	if err != nil {
		logger.Error("iam", "register failed",
			zap.String("username", cmd.Username),
			zap.Error(err),
		)

		// 根据错误类型返回不同的状态码
		switch err {
		case domain.ErrUserAlreadyExists:
			response.Error(c, http.StatusConflict, "用户名已存在")
		default:
			response.Error(c, http.StatusInternalServerError, "注册失败："+err.Error())
		}
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
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.appService.Login(platform_http.Ctx(c), cmd, h.tokenManager)
	if err != nil {
		logger.Warn("iam", "login failed",
			zap.String("username", cmd.Username),
			zap.Error(err),
		)

		// 根据错误类型返回不同的状态码
		switch err {
		case domain.ErrUserNotFound:
			response.Error(c, http.StatusUnauthorized, "用户不存在")
		case domain.ErrUserFrozen:
			response.Error(c, http.StatusForbidden, "账号已被冻结")
		case domain.ErrInvalidPassword:
			response.Error(c, http.StatusUnauthorized, "密码错误")
		default:
			response.Error(c, http.StatusInternalServerError, "登录失败："+err.Error())
		}
		return
	}

	logger.Info("iam", "user logged in successfully",
		zap.String("username", cmd.Username),
		zap.Int64("user_id", resp.UserInfo.ID),
	)
	response.Success(c, resp)
}
