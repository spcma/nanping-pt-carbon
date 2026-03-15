package http

import (
	"app/internal/module/iam/application"
	"app/internal/module/iam/domain"
	http2 "app/internal/platform/http"
	"app/internal/platform/http/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SysRoleHandler system role handler
type SysRoleHandler struct {
	appService *application.SysRoleAppService
}

// NewSysRoleHandler creates system role handler
func NewSysRoleHandler(appService *application.SysRoleAppService) *SysRoleHandler {
	return &SysRoleHandler{
		appService: appService,
	}
}

// Create creates a system role
func (h *SysRoleHandler) Create(c *gin.Context) {
	var cmd application.CreateSysRoleCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user := http2.GetCurrentUser(c)
	cmd.UserID = user.ID

	id, err := h.appService.CreateSysRole(http2.Ctx(c), cmd)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"id": id})
}

// Update updates a system role
func (h *SysRoleHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	var cmd application.UpdateSysRoleCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	cmd.ID = id

	user := http2.GetCurrentUser(c)
	cmd.UserID = user.ID

	if err := h.appService.UpdateSysRole(http2.Ctx(c), cmd); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete deletes a system role
func (h *SysRoleHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.appService.DeleteSysRole(http2.Ctx(c), id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetByID gets system role by ID
func (h *SysRoleHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}

	role, err := h.appService.GetSysRoleByID(http2.Ctx(c), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, role)
}

// GetPage queries system roles with pagination
func (h *SysRoleHandler) GetPage(c *gin.Context) {
	var query domain.SysRolePageQuery
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

	roles, total, err := h.appService.GetSysRolePage(http2.Ctx(c), query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":  roles,
		"total": total,
	})
}

// ChangeStatus changes role status
func (h *SysRoleHandler) ChangeStatus(c *gin.Context) {
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

	user := http2.GetCurrentUser(c)
	if err := h.appService.ChangeRoleStatus(http2.Ctx(c), application.ChangeRoleStatusCommand{ID: id, Status: domain.RoleStatus(cmd.Status), UserID: user.ID}); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}
