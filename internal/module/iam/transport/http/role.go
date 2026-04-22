package http

import (
	"app/internal/module/iam/application"
	"app/internal/module/iam/domain"
	platform_http "app/internal/platform/http"
	"app/internal/platform/http/response"
	"app/internal/shared/logger"
	"net/http"
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
	var cmd application.CreateSysRoleCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if user := platform_http.GetCurrentUser(c); user != nil {
		cmd.UserID = user.ID
	}

	id, err := h.appService.CreateRole(platform_http.Ctx(c), cmd)
	if err != nil {
		logger.RuntimeL.WithTraceID(platform_http.GetTraceID(c)).Error("角色创建失败", zap.Error(err))
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"id": id})
}

// Update updates a system role
func (h *RoleHandler) Update(c *gin.Context) {

	var cmd application.UpdateSysRoleCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if cmd.ID == 0 {
		response.Error(c, http.StatusBadRequest, "id is required")
		return
	}

	if user := platform_http.GetCurrentUser(c); user != nil {
		cmd.UserID = user.ID
	}

	if err := h.appService.UpdateRole(platform_http.Ctx(c), cmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, "修改成功")
}

// Delete deletes a system role
func (h *RoleHandler) Delete(c *gin.Context) {
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

	if user := platform_http.GetCurrentUser(c); user != nil {
		cmd.UserID = user.ID
	}

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
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	role, err := h.appService.GetRoleByID(platform_http.Ctx(c), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, role)
}

// GetPage queries system roles with pagination
func (h *RoleHandler) GetPage(c *gin.Context) {
	var query domain.SysRolePageQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
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
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, res)
}

// ChangeStatus changes role status
func (h *RoleHandler) ChangeStatus(c *gin.Context) {
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
	if err := h.appService.ChangeRoleStatus(platform_http.Ctx(c), application.ChangeRoleStatusCommand{ID: id, Status: domain.RoleStatus(cmd.Status), UserID: user.ID}); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
