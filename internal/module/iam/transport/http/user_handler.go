package http

import (
	"app/internal/module/iam/application"
	"app/internal/module/iam/domain"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"

	"github.com/spf13/cast"
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
	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "")
		return
	}

	var cmd application.CreateUserCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "")
		return
	}

	cmd.UserID = currentUser.ID

	id, err := h.appService.Create(platform_http.Ctx(c), cmd)
	if err != nil {
		logger.IamL.Error("create user failed",
			zap.String("username", *cmd.Username),
			zap.Error(err),
		)
		response.InternalError(c, "创建用户失败")
		return
	}

	logger.IamL.Info("currentUser created successfully",
		zap.Int64("user_id", id),
		zap.String("username", *cmd.Username),
		zap.Int64("create_by", currentUser.ID))

	response.Success(c, gin.H{"id": id})
}

// Update updates a system user
func (h *UserHandler) Update(c *gin.Context) {
	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "")
		return
	}

	var cmd application.UpdateUserCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	cmd.UserID = currentUser.ID

	if err := h.appService.Update(platform_http.Ctx(c), cmd); err != nil {
		logger.IamL.Error("update user failed",
			zap.Int64("user_id", cmd.ID),
			zap.Int64("update_by", currentUser.ID),
			zap.Error(err),
		)
		response.InternalError(c, "更新用户信息失败")
		return
	}

	response.Success(c, nil)
}

// Delete deletes a system user
func (h *UserHandler) Delete(c *gin.Context) {
	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "")
		return
	}

	var cmd application.DeleteUserCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if cmd.ID == 0 {
		response.BadRequest(c, "无效的ID")
		return
	}

	cmd.UserID = currentUser.ID

	if err := h.appService.Delete(platform_http.Ctx(c), &cmd); err != nil {
		logger.IamL.Error("delete user failed",
			zap.Int64("user_id", cmd.ID),
			zap.Int64("delete_by", currentUser.ID),
			zap.Error(err),
		)

		response.InternalError(c, "删除用户失败")
		return
	}

	response.Success(c, nil)
}

func (h *UserHandler) GetById(c *gin.Context) {
	idStr := c.Query("id")
	id := cast.ToInt64(idStr)
	if id == 0 {
		response.BadRequest(c, "无效的ID")
		return
	}

	user, err := h.appService.GetByID(platform_http.Ctx(c), id)
	if err != nil {
		logger.IamL.Error("get user failed",
			zap.Int64("user_id", id),
			zap.Error(err),
		)
		response.InternalError(c, "获取用户信息失败")
		return
	}

	response.Success(c, user)
}

// GetByQuery
func (h *UserHandler) GetByQuery(c *gin.Context) {
	var query domain.UserQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "")
		return
	}

	user, err := h.appService.GetByQuery(platform_http.Ctx(c), &query)
	if err != nil {
		logger.IamL.Error("get user by query failed",
			zap.Error(err),
		)
		response.InternalError(c, "查询用户失败")
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) GetList(c *gin.Context) {
	list, err := h.appService.GetList(platform_http.Ctx(c))
	if err != nil {
		logger.IamL.Error("get user list failed",
			zap.Error(err),
		)
		response.InternalError(c, "获取用户列表失败")
		return
	}

	response.Success(c, list)
}

// GetPage queries system users with pagination（需要认证）
func (h *UserHandler) GetPage(c *gin.Context) {
	var query domain.UsersPageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "")
		return
	}

	query.Fixed()

	res, err := h.appService.GetPage(platform_http.Ctx(c), &query)
	if err != nil {
		logger.IamL.Error("get user page failed",
			zap.Int("page_num", query.PageNum),
			zap.Int("page_size", query.PageSize),
			zap.Error(err),
		)
		response.InternalError(c, "分页查询用户失败")
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
		response.BadRequest(c, "请求参数错误")
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
		logger.IamL.Error("get public user page failed",
			zap.Int("page_num", query.PageNum),
			zap.Int("page_size", query.PageSize),
			zap.Error(err),
		)
		response.InternalError(c, "分页查询用户失败")
		return
	}

	// 获取当前用户信息（可能为 nil）
	currentUser := platform_http.GetCurrentUser(c)
	// 根据是否登录返回不同信息
	if currentUser != nil {
		// 已登录：返回增强信息
		logger.IamL.Info("public page accessed by authenticated user", zap.Int64("user_id", currentUser.ID))

		h.respondWithEnhancedUsers(c, res.Data, res.Total)
	} else {
		// 未登录：返回基础公开信息
		logger.IamL.Info("public page accessed by anonymous user")

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
		response.BadRequest(c, "")
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
		logger.IamL.Error("reset password failed",
			zap.Int64("user_id", cmd.Id),
			zap.Error(err))

		response.InternalError(c, "修改密码失败")
		return
	}

	response.Success(c, "重置成功")
}

// ChangePassword changes user password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var cmd application.ChangePasswordCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "")
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
		logger.IamL.Error("change password failed",
			zap.Int64("user_id", cmd.Id),
			zap.Error(err),
		)

		response.InternalError(c, "修改密码失败")
		return
	}

	response.Success(c, "密码修改成功")
}

// ChangeStatus changes user status
func (h *UserHandler) ChangeStatus(c *gin.Context) {

	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "")
		return
	}

	type changStatusRequest struct {
		ID     int64  `json:"id"`
		Status string `json:"status"`
	}

	var cmd changStatusRequest
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "")
		return
	}

	if cmd.ID == 0 {
		response.BadRequest(c, "id is required")
		return
	}

	if err := h.appService.ChangeUserStatus(platform_http.Ctx(c), cmd.ID, domain.UserStatus(cmd.Status), currentUser.ID); err != nil {
		logger.IamL.Error("change user status failed",
			zap.Int64("user_id", cmd.ID),
			zap.String("status", cmd.Status),
			zap.Error(err),
		)
		response.InternalError(c, "变更用户状态失败")
		return
	}

	response.Success(c, "状态修改成功")
}
