package http

import (
	"app/internal/module/iam/application"
	"app/internal/module/iam/domain"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RoleHandler system role handler
type RoleHandler struct {
	appService *application.RoleAppService
}

// NewRoleHandler creates system role handler
func NewRoleHandler(appService *application.RoleAppService) *RoleHandler {
	return &RoleHandler{
		appService: appService,
	}
}

// Create creates a system role
func (h *RoleHandler) Create(c *gin.Context) {

	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "")
		return
	}

	var cmd application.CreateSysRoleCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cmd.UserID = currentUser.ID

	id, err := h.appService.CreateRole(platform_http.Ctx(c), cmd)
	if err != nil {
		logger.RuntimeL.Error("create role failed",
			zap.String("role_name", cmd.Name),
			zap.Error(err),
		)
		response.InternalError(c, "创建角色失败")
		return
	}

	response.Success(c, gin.H{"id": id})
}

// Update updates a system role
func (h *RoleHandler) Update(c *gin.Context) {
	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "")
		return
	}

	var cmd application.UpdateSysRoleCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if cmd.ID == 0 {
		response.BadRequest(c, "ID不能为空")
		return
	}

	cmd.UserID = currentUser.ID

	if err := h.appService.UpdateRole(platform_http.Ctx(c), cmd); err != nil {
		logger.RuntimeL.Error("update role failed",
			zap.Int64("role_id", cmd.ID),
			zap.Error(err),
		)
		response.InternalError(c, "更新角色失败")
		return
	}

	response.Success(c, "修改成功")
}

// Delete deletes a system role
func (h *RoleHandler) Delete(c *gin.Context) {
	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "")
		return
	}

	type DeleteSysRoleCommand struct {
		ID     int64 `json:"id"`
		UserID int64 `json:"userid"`
	}

	var cmd DeleteSysRoleCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "invalid request")
		return
	}

	if cmd.ID == 0 {
		response.BadRequest(c, "id is required")
		return
	}

	cmd.UserID = currentUser.ID

	if cmd.UserID == 0 {
		response.BadRequest(c, "user is required")
		return
	}

	if err := h.appService.DeleteRole(platform_http.Ctx(c), cmd.ID, cmd.UserID); err != nil {
		logger.RuntimeL.WithTraceID(platform_http.GetTraceID(c)).Error("角色删除失败", zap.Error(err))
		response.InternalError(c, "删除失败")
		return
	}

	response.Success(c, "删除成功")
}

// GetByID gets system role by ID
func (h *RoleHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	role, err := h.appService.GetRoleByID(platform_http.Ctx(c), id)
	if err != nil {
		logger.RuntimeL.Error("get role failed",
			zap.Int64("role_id", id),
			zap.Error(err),
		)
		response.InternalError(c, "获取角色失败")
		return
	}

	response.Success(c, role)
}

// GetPage queries system roles with pagination
func (h *RoleHandler) GetPage(c *gin.Context) {
	var query domain.SysRolePageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	query.Fixed()

	if query.PageNum == 0 {
		query.PageNum = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 10
	}

	res, err := h.appService.GetRolePage(platform_http.Ctx(c), &query)
	if err != nil {
		logger.RuntimeL.Error("get role page failed",
			zap.Int("page_num", query.PageNum),
			zap.Int("page_size", query.PageSize),
			zap.Error(err),
		)
		response.InternalError(c, "分页查询角色失败")
		return
	}

	response.Success(c, res)
}

// ChangeStatus changes role status
func (h *RoleHandler) ChangeStatus(c *gin.Context) {

	currentUser := platform_http.GetCurrentUser(c)
	if currentUser == nil {
		response.Unauthorized(c, "")
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var cmd struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	if err = h.appService.ChangeRoleStatus(platform_http.Ctx(c), application.ChangeRoleStatusCommand{ID: id, Status: domain.RoleStatus(cmd.Status), UserID: currentUser.ID}); err != nil {
		logger.RuntimeL.Error("change role status failed",
			zap.Int64("role_id", id),
			zap.String("status", cmd.Status),
			zap.Error(err),
		)
		response.InternalError(c, "变更角色状态失败")
		return
	}

	response.Success(c, nil)
}
