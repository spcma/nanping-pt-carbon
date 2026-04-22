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

// UserHandler system user handler
type UserHandler struct {
	appService *application.UserService
}

// NewUserHandler creates system user handler
func NewUserHandler(appService *application.UserService) *UserHandler {
	return &UserHandler{
		appService: appService,
	}
}

// Create creates a system user
func (h *UserHandler) Create(c *gin.Context) {
	var cmd application.CreateUserCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		logger.Warn("iam", "create currentUser - invalid request",
			zap.String("error", err.Error()),
		)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	//currentUser := platform_http.GetCurrentUser(c)
	//if currentUser == nil {
	//	response.BadRequest(c, "user not found")
	//	return
	//}
	//cmd.UserID = currentUser.ID

	id, err := h.appService.Create(platform_http.Ctx(c), cmd)
	if err != nil {
		logger.Error("iam", "create currentUser failed",
			zap.String("username", *cmd.Username),
			zap.Error(err),
		)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("iam", "currentUser created successfully",
		zap.Int64("user_id", id),
		zap.String("username", *cmd.Username),
	)
	response.Success(c, gin.H{"id": id})
}

// Update updates a system user
func (h *UserHandler) Update(c *gin.Context) {
	var cmd application.UpdateUserCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.appService.Update(platform_http.Ctx(c), cmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete deletes a system user
func (h *UserHandler) Delete(c *gin.Context) {
	var param application.DeleteUserCommand
	if err := c.ShouldBindJSON(&param); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user := platform_http.GetCurrentUser(c)

	param.UserID = user.ID

	if err := h.appService.Delete(platform_http.Ctx(c), &param); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *UserHandler) GetById(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.appService.GetByID(platform_http.Ctx(c), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, user)
}

// GetByQuery
func (h *UserHandler) GetByQuery(c *gin.Context) {
	var query domain.UserQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.appService.GetByQuery(platform_http.Ctx(c), &query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) GetList(c *gin.Context) {
	list, err := h.appService.GetList(platform_http.Ctx(c))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, list)
}

// GetPage queries system users with pagination（需要认证）
func (h *UserHandler) GetPage(c *gin.Context) {
	var query domain.UsersPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	query.Fixed()

	res, err := h.appService.GetPage(platform_http.Ctx(c), &query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, res)
}

// GetPublicPage queries public system users with pagination（支持可选认证）
// - 未登录：返回基础公开信息
// - 已登录：返回增强信息（包含更多字段）
func (h *UserHandler) GetPublicPage(c *gin.Context) {
	var query domain.UsersPageQuery
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

	res, err := h.appService.GetPage(platform_http.Ctx(c), &query)
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
func (h *UserHandler) respondWithPublicUsers(c *gin.Context, users []*domain.User, total int64) {
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
func (h *UserHandler) respondWithEnhancedUsers(c *gin.Context, users []*domain.User, total int64) {
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

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var cmd application.ChangePasswordCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if cmd.Id == 0 {
		response.BadRequest(c, "id is required")
		return
	}

	if cmd.NewPassword == "" {
		response.BadRequest(c, "newPassword is required")
		return
	}

	//	仅管理员可以强制修改密码，后续需要加 force 赋值校验
	err := h.appService.ChangePassword(platform_http.Ctx(c), &cmd)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, "重置成功")
}

// ChangePassword changes user password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var cmd application.ChangePasswordCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if cmd.Id == 0 {
		response.BadRequest(c, "id is required")
		return
	}

	if cmd.NewPassword == "" {
		response.BadRequest(c, "newPassword is required")
		return
	}

	err := h.appService.ChangePassword(platform_http.Ctx(c), &cmd)
	if err != nil {
		response.InternalError(c, "change password failed")
		return
	}

	response.Success(c, "密码修改成功")
}

// ChangeStatus changes user status
func (h *UserHandler) ChangeStatus(c *gin.Context) {
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
