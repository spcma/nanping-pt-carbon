package http

import (
	"app/internal/module/iam/application"
	"app/internal/module/iam/domain"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
	"net/http"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// SysUserHandler system user handler
type SysUserHandler struct {
	appService *application.UsersService
}

// NewSysUserHandler creates system user handler
func NewSysUserHandler(appService *application.UsersService) *SysUserHandler {
	return &SysUserHandler{
		appService: appService,
	}
}

// Create creates a system user
func (h *SysUserHandler) Create(c *gin.Context) {
	var cmd application.CreateSysUserCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		logger.Warn("iam", "create securityUser - invalid request",
			zap.String("error", err.Error()),
		)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	securityUser := platform_http.GetCurrentUser(c)
	cmd.UserID = securityUser.ID

	id, err := h.appService.Create(platform_http.Ctx(c), cmd)
	if err != nil {
		logger.Error("iam", "create securityUser failed",
			zap.String("username", cmd.Username),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("iam", "securityUser created successfully",
		zap.Int64("user_id", id),
		zap.String("username", cmd.Username),
	)
	response.Success(c, gin.H{"id": id})
}

// Update updates a system user
func (h *SysUserHandler) Update(c *gin.Context) {
	var cmd application.UpdateSysUserParam
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user := platform_http.GetCurrentUser(c)
	cmd.UserID = user.ID

	if err := h.appService.UpdateSysUser(platform_http.Ctx(c), cmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete deletes a system user
func (h *SysUserHandler) Delete(c *gin.Context) {
	var param application.DeleteSysUserParam
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user := platform_http.GetCurrentUser(c)

	param.UserID = user.ID

	if err := h.appService.DeleteSysUser(platform_http.Ctx(c), &param); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *SysUserHandler) GetById(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.appService.GetSysUserByID(platform_http.Ctx(c), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, user)
}

// GetByQuery
func (h *SysUserHandler) GetByQuery(c *gin.Context) {
	var query application.UsersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

}

func (h *SysUserHandler) GetList(c *gin.Context) {
	list, err := h.appService.GetList(platform_http.Ctx(c))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, list)
}

// GetPage queries system users with pagination（需要认证）
func (h *SysUserHandler) GetPage(c *gin.Context) {
	var query domain.SysUserPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	query.Fixed()

	res, err := h.appService.GetSysUserPage(platform_http.Ctx(c), &query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, res)
}

// GetPublicPage queries public system users with pagination（支持可选认证）
// - 未登录：返回基础公开信息
// - 已登录：返回增强信息（包含更多字段）
func (h *SysUserHandler) GetPublicPage(c *gin.Context) {
	var query domain.SysUserPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if query.PageNum == 0 {
		query.PageNum = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	res, err := h.appService.GetSysUserPage(platform_http.Ctx(c), &query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	// 获取当前用户信息（可能为 nil）
	securityUser := platform_http.GetCurrentUser(c)

	// 根据是否登录返回不同信息
	if securityUser != nil {
		// 已登录：返回增强信息
		logger.Debug("iam", "public page accessed by authenticated user",
			zap.Int64("user_id", securityUser.ID),
		)
		h.respondWithEnhancedUsers(c, res.Data, res.Total)
	} else {
		// 未登录：返回基础公开信息
		logger.Debug("iam", "public page accessed by anonymous user")
		h.respondWithPublicUsers(c, res.Data, res.Total)
	}
}

// respondWithPublicUsers 返回基础公开信息（脱敏）
func (h *SysUserHandler) respondWithPublicUsers(c *gin.Context, users []*domain.Users, total int64) {
	type PublicUserInfo struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}

	publicUsers := make([]*PublicUserInfo, 0, len(users))
	for _, u := range users {
		publicUsers = append(publicUsers, &PublicUserInfo{
			ID:       u.Id,
			Username: u.Username,
			Nickname: u.Nickname,
			Avatar:   u.Avatar,
		})
	}

	response.Success(c, gin.H{
		"list":  publicUsers,
		"total": total,
	})
}

// respondWithEnhancedUsers 返回增强信息（包含更多字段）
func (h *SysUserHandler) respondWithEnhancedUsers(c *gin.Context, users []*domain.Users, total int64) {
	type EnhancedUserInfo struct {
		ID         int64  `json:"id"`
		Username   string `json:"username"`
		Nickname   string `json:"nickname"`
		Avatar     string `json:"avatar"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		CreateTime string `json:"createTime"`
	}

	enhancedUsers := make([]*EnhancedUserInfo, 0, len(users))
	for _, u := range users {
		enhancedUsers = append(enhancedUsers, &EnhancedUserInfo{
			ID:         u.Id,
			Username:   u.Username,
			Nickname:   u.Nickname,
			Avatar:     u.Avatar,
			Email:      u.Email,
			Phone:      u.Phone,
			CreateTime: u.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	response.Success(c, gin.H{
		"list":  enhancedUsers,
		"total": total,
	})
}

// ChangePassword changes user password
func (h *SysUserHandler) ChangePassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var cmd application.ChangePasswordCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	cmd.Id = id

	// ChangePasswordCommand 不需要 UserID，只需要 Id 和 Password

	if err := h.appService.ChangePassword(platform_http.Ctx(c), cmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// ChangeStatus changes user status
func (h *SysUserHandler) ChangeStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var cmd struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user := platform_http.GetCurrentUser(c)
	if err := h.appService.ChangeUserStatus(platform_http.Ctx(c), id, domain.UserStatus(cmd.Status), user.ID); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
